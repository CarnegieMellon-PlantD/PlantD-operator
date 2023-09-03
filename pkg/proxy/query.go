package proxy

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/CarnegieMellon-PlantD/PlantD-operator/pkg/config"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var (
	step time.Duration
)

func init() {
	scrapeInterval := config.GetString("database.prometheus.scrapeInterval")
	interval, err := time.ParseDuration(scrapeInterval)
	if err != nil {
		panic(err)
	}
	step = 2 * interval
}

type QueryClient interface {
	QueryBiChann(ctx context.Context, req *QueryRequest) (*QueryResult, error)
	QueryTriChann(ctx context.Context, req *QueryRequest) (*QueryResult, error)
}

type QueryAgent struct {
	PromAPI v1.API
}

type LabelSelector struct {
	TargetLabels []string
}

func NewLabelSelector(labels []string) *LabelSelector {
	return &LabelSelector{
		TargetLabels: labels,
	}
}

func (s *LabelSelector) GetSeriesFromMetric(labelSet model.Metric) string {
	values := make([]string, len(s.TargetLabels))
	for i, v := range s.TargetLabels {
		values[i] = string(labelSet[model.LabelName(v)])
	}
	return strings.Join(values, "/")
}

func NewQueryAgent(url string) (*QueryAgent, error) {
	promClient, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		return nil, err
	}
	v1api := v1.NewAPI(promClient)
	return &QueryAgent{
		PromAPI: v1api,
	}, nil
}

func (c *QueryAgent) QueryBiChann(ctx context.Context, req *QueryRequest) (*QueryResult, error) {
	switch req.Source {
	case DATA_SOURCE_PROM:
		return c.QueryBiChannProm(ctx, req)
	case DATA_SOURCE_REDIS:
		return c.QueryBiChannRedis(ctx, req)
	}

	return nil, fmt.Errorf("data source: %s not supported", req.Source)
}

func (c *QueryAgent) QueryTriChann(ctx context.Context, req *QueryRequest) (*QueryResult, error) {
	switch req.Source {
	case DATA_SOURCE_PROM:
		return c.QueryTriChannProm(ctx, req)
	case DATA_SOURCE_REDIS:
		return c.QueryTriChannRedis(ctx, req)
	}

	return nil, fmt.Errorf("data source: %s not supported", req.Source)
}

func SampleValueToFloatPointer(v model.SampleValue) *float64 {
	value := float64(v)
	if math.IsNaN(value) {
		return nil
	}
	return &value
}

func (c *QueryAgent) QueryBiChannProm(ctx context.Context, req *QueryRequest) (*QueryResult, error) {
	result, _, err := c.PromAPI.Query(ctx, req.Param.Query, req.Param.EndTimestamp.Time)
	if err != nil {
		return nil, err
	}
	labelSelector := NewLabelSelector(req.Param.LabelSelector)
	if vectorVal, ok := result.(model.Vector); ok {
		n := len(vectorVal)
		res := make([]ResultPoint, n)
		for i, item := range vectorVal {
			res[i] = ResultPoint{
				Series: labelSelector.GetSeriesFromMetric(item.Metric),
				ValueY: SampleValueToFloatPointer(item.Value),
			}
		}
		return &QueryResult{
			Result: res,
		}, nil
	}
	return nil, fmt.Errorf("cannot convert data to desired format")
}

func (c *QueryAgent) QueryTriChannProm(ctx context.Context, req *QueryRequest) (*QueryResult, error) {
	promStep := step
	if req.Param.Step > 0 {
		promStep = time.Duration(req.Param.Step) * time.Second
	}

	r := v1.Range{
		Start: req.Param.StartTimestamp.Time,
		End:   req.Param.EndTimestamp.Time,
		Step:  promStep,
	}
	result, _, err := c.PromAPI.QueryRange(ctx, req.Param.Query, r)
	if err != nil {
		return nil, err
	}
	labelSelector := NewLabelSelector(req.Param.LabelSelector)
	if matrixVal, ok := result.(model.Matrix); ok {
		res := []ResultPoint{}
		for _, series := range matrixVal {
			s := labelSelector.GetSeriesFromMetric(series.Metric)
			for _, point := range series.Values {
				ts := float64(point.Timestamp.Time().Unix())
				res = append(res, ResultPoint{
					Series: s,
					ValueY: SampleValueToFloatPointer(point.Value),
					ValueX: &ts,
				})
			}
		}
		return &QueryResult{
			Result: res,
		}, nil
	}
	return nil, fmt.Errorf("cannot convert data to desired format")
}

func (c *QueryAgent) QueryBiChannRedis(ctx context.Context, req *QueryRequest) (*QueryResult, error) {

	return nil, nil
}

func (c *QueryAgent) QueryTriChannRedis(ctx context.Context, req *QueryRequest) (*QueryResult, error) {

	return nil, nil
}
