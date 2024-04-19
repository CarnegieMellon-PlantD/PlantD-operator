package proxy

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"

	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/redis/go-redis/v9"
)

var (
	promUrl   string
	redisAddr string
)

func init() {
	promUrl = config.GetString("database.prometheus.thanosUrl")
	redisAddr = fmt.Sprintf("%s:%d", config.GetString("database.redis.host"), config.GetInt("database.redis.port"))
}

type QueryAgent struct {
	PromAPI       prometheusv1.API
	RedisClient   *redis.Client
	RedisTSClient *redistimeseries.Client
}

func NewQueryAgent() (*QueryAgent, error) {
	promClient, err := api.NewClient(api.Config{
		Address: promUrl,
	})
	if err != nil {
		return nil, err
	}
	promApi := prometheusv1.NewAPI(promClient)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	redisTSClient := redistimeseries.NewClient(redisAddr, "redis-ts-client", nil)

	return &QueryAgent{
		PromAPI:       promApi,
		RedisClient:   redisClient,
		RedisTSClient: redisTSClient,
	}, nil
}

type LabelSelector struct {
	TargetLabels []string
	pastSeries   map[string]struct{}
}

func NewLabelSelector(labels []string) *LabelSelector {
	return &LabelSelector{
		TargetLabels: labels,
		pastSeries:   make(map[string]struct{}),
	}
}

func (ls *LabelSelector) GetSeriesFromPromMetric(m model.Metric) (string, error) {
	labelValues := make([]string, len(ls.TargetLabels))
	for i, v := range ls.TargetLabels {
		if labelValue, ok := m[model.LabelName(v)]; ok {
			labelValues[i] = string(labelValue)
		}
	}
	result := strings.Join(labelValues, "/")
	if _, ok := ls.pastSeries[result]; ok {
		return "", fmt.Errorf("duplicated series")
	}
	ls.pastSeries[result] = struct{}{}
	return result, nil
}

func (ls *LabelSelector) GetSeriesFromRedisRange(r redistimeseries.Range) (string, error) {
	labelValues := make([]string, len(ls.TargetLabels))
	for i, v := range ls.TargetLabels {
		if labelValue, ok := r.Labels[v]; ok {
			labelValues[i] = labelValue
		}
	}
	result := strings.Join(labelValues, "/")
	if _, ok := ls.pastSeries[result]; ok {
		return "", fmt.Errorf("duplicated series")
	}
	ls.pastSeries[result] = struct{}{}
	return result, nil
}

func getFloat64PtrFromFloat64(n float64) *float64 {
	if math.IsNaN(n) {
		return nil
	}
	return &n
}

func (qa *QueryAgent) PromQuery(ctx context.Context, req *PromRequest) (*BiChanResponse, error) {
	result, _, err := qa.PromAPI.Query(ctx, req.Query, req.EndTimestamp.Time)
	if err != nil {
		return nil, err
	}
	labelSelector := NewLabelSelector(req.LabelSelector)
	if vectorVal, ok := result.(model.Vector); ok {
		res := make([]*BiChanDataPoint, len(vectorVal))
		for i, sampleVal := range vectorVal {
			series, err := labelSelector.GetSeriesFromPromMetric(sampleVal.Metric)
			if err != nil {
				return nil, err
			}
			res[i] = &BiChanDataPoint{
				Series: series,
				ValueY: getFloat64PtrFromFloat64(float64(sampleVal.Value)),
			}
		}
		return &BiChanResponse{
			Result: res,
		}, nil
	}
	return nil, fmt.Errorf("cannot convert data to desired format")
}

func (qa *QueryAgent) PromQueryRange(ctx context.Context, req *PromRequest) (*TriChanResponse, error) {
	timeRange := prometheusv1.Range{
		Start: req.StartTimestamp.Time,
		End:   req.EndTimestamp.Time,
		Step:  time.Duration(req.Step) * time.Second,
	}
	result, _, err := qa.PromAPI.QueryRange(ctx, req.Query, timeRange)
	if err != nil {
		return nil, err
	}
	labelSelector := NewLabelSelector(req.LabelSelector)
	if matrixVal, ok := result.(model.Matrix); ok {
		res := make([]*TriChanDataPoint, 0)
		for _, streamVal := range matrixVal {
			series, err := labelSelector.GetSeriesFromPromMetric(streamVal.Metric)
			if err != nil {
				return nil, err
			}
			for _, sampleVal := range streamVal.Values {
				res = append(res, &TriChanDataPoint{
					Series: series,
					ValueY: getFloat64PtrFromFloat64(float64(sampleVal.Value)),
					ValueX: getFloat64PtrFromFloat64(float64(sampleVal.Timestamp.Time().Unix())),
				})
			}
		}
		return &TriChanResponse{
			Result: res,
		}, nil
	}
	return nil, fmt.Errorf("cannot convert data to desired format")
}

func (qa *QueryAgent) RedisGet(ctx context.Context, req *RedisRequest) (*RawResponse, error) {
	val, err := qa.RedisClient.Get(ctx, req.Key).Result()
	if err != nil {
		return nil, err
	}
	return &RawResponse{
		Result: val,
	}, nil
}

func (qa *QueryAgent) RedisTSMultiGet(ctx context.Context, req *RedisTSRequest) (*BiChanResponse, error) {
	resultCh := make(chan *BiChanResponse)
	errorCh := make(chan error)
	go func() {
		opts := redistimeseries.NewMultiGetOptions()
		opts.WithLabels = true
		rangesVal, err := qa.RedisTSClient.MultiGetWithOptions(*opts, req.Filters...)
		if err != nil {
			errorCh <- err
			return
		}
		labelSelector := NewLabelSelector(req.LabelSelector)
		res := make([]*BiChanDataPoint, len(rangesVal))
		for i, rangeVal := range rangesVal {
			series, err := labelSelector.GetSeriesFromRedisRange(rangeVal)
			if err != nil {
				errorCh <- err
				return
			}
			if len(rangeVal.DataPoints) != 1 {
				errorCh <- fmt.Errorf("expect 1 data point per range but got %d", len(rangeVal.DataPoints))
				return
			}
			res[i] = &BiChanDataPoint{
				Series: series,
				ValueY: getFloat64PtrFromFloat64(rangeVal.DataPoints[0].Value),
			}
		}
		resultCh <- &BiChanResponse{
			Result: res,
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return result, nil
	case err := <-errorCh:
		return nil, err
	}
}

func (qa *QueryAgent) RedisTSMultiRange(ctx context.Context, req *RedisTSRequest) (*TriChanResponse, error) {
	resultCh := make(chan *TriChanResponse)
	errorCh := make(chan error)
	go func() {
		opts := redistimeseries.NewMultiRangeOptions()
		opts.WithLabels = true
		rangesVal, err := qa.RedisTSClient.MultiRangeWithOptions(req.StartTimestamp.UnixMilli(), req.EndTimestamp.UnixMilli(), *opts, req.Filters...)
		if err != nil {
			errorCh <- err
			return
		}
		labelSelector := NewLabelSelector(req.LabelSelector)
		res := make([]*TriChanDataPoint, 0)
		for _, rangeVal := range rangesVal {
			series, err := labelSelector.GetSeriesFromRedisRange(rangeVal)
			if err != nil {
				errorCh <- err
				return
			}
			for _, dataPointVal := range rangeVal.DataPoints {
				res = append(res, &TriChanDataPoint{
					Series: series,
					ValueY: getFloat64PtrFromFloat64(dataPointVal.Value),
					// Redis uses Unix timestamp in milliseconds, convert it to Unix timestamp in seconds
					ValueX: getFloat64PtrFromFloat64(float64(dataPointVal.Timestamp) / 1000),
				})
			}
		}
		resultCh <- &TriChanResponse{
			Result: res,
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultCh:
		return result, nil
	case err := <-errorCh:
		return nil, err
	}
}
