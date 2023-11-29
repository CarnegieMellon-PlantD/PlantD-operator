package proxy

import (
	"encoding/json"
	"time"

	"k8s.io/apimachinery/pkg/types"
)

// ErrorResponse defines the response to send when error occurs.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// CheckHTTPHealthRequest defines the request to check health of a URL using HTTP protocol.
type CheckHTTPHealthRequest struct {
	URL string `json:"url,omitempty"`
}

// ImportResourcesStatistics contains the result of importing resources.
type ImportResourcesStatistics struct {
	// NumSucceeded is the number of resources that are successfully imported
	NumSucceeded int `json:"numSucceeded"`
	// NumFailed is the number of resources that failed to be imported
	NumFailed int `json:"numFailed"`
	// ErrorMessages contains the error messages if any
	ErrorMessages []string `json:"errors"`
}

// ExportResourceInfo contains the kind, namespace, and name that are necessary to locate a unique resource to
// export. Note that the group and version are fixed across all resources and are thus omitted.
type ExportResourceInfo struct {
	Kind string
	types.NamespacedName
}

// SourceType is the type of data source
type SourceType int8

const (
	Prometheus SourceType = iota
	RedisTimeSeries
)

// ChanType is the type of channel number per data point in response data
type ChanType int8

const (
	BiChan ChanType = iota
	TriChan
)

// UnixTimestamp is the Unix timestamp in second.
type UnixTimestamp struct {
	time.Time
}

func (ts *UnixTimestamp) UnmarshalJSON(b []byte) error {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	ts.Time = time.Unix(timestamp, 0)
	return nil
}

// PromRequest contains the parameters for making a "Query" or "QueryRange" request to Prometheus.
type PromRequest struct {
	Query          string        `json:"query,omitempty"`
	StartTimestamp UnixTimestamp `json:"start,omitempty"`
	EndTimestamp   UnixTimestamp `json:"end,omitempty"`
	Step           int64         `json:"step,omitempty"`
	LabelSelector  []string      `json:"labelSelector,omitempty"`
}

// RedisTSRequest contains the parameters for making a "MultiGet" or "MultiRange" request to Redis Time Series.
type RedisTSRequest struct {
	Filters        []string      `json:"filters,omitempty"`
	StartTimestamp UnixTimestamp `json:"start,omitempty"`
	EndTimestamp   UnixTimestamp `json:"end,omitempty"`
	LabelSelector  []string      `json:"labelSelector,omitempty"`
}

// BiChanDataPoint defines the data point in bi-channel data.
type BiChanDataPoint struct {
	Series string  `json:"series"`
	ValueY float64 `json:"y"`
}

// BiChanResponse defines the response to send for bi-channel data.
type BiChanResponse struct {
	Result []*BiChanDataPoint `json:"result"`
}

// TriChanDataPoint defines the data point in tri-channel data.
type TriChanDataPoint struct {
	Series string  `json:"series"`
	ValueY float64 `json:"y"`
	ValueX float64 `json:"x"`
}

// TriChanResponse defines the response to send for tri-channel data.
type TriChanResponse struct {
	Result []*TriChanDataPoint `json:"result"`
}
