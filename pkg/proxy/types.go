package proxy

import (
	"encoding/json"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NamespaceInfo represents information about a namespace.
type NamespaceInfo struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// CustomResourceMeta represents metadata information for a custom resource.
type CustomResourceMeta struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Kind      string `json:"kind,omitempty"`
}

// ExportResourcesInfo represents information about custom resources to be exported.
type ExportResourcesInfo struct {
	Metadata []CustomResourceMeta `json:"metadata,omitempty"`
}

// CustomResourceWrapper is a wrapper for a custom resource that implements the client.Object interface.
type CustomResourceWrapper struct {
	client.Object
}

// HealthCheckMeta represents metadata for health check information.
type HealthCheckMeta struct {
	URL                 string `json:"url,omitempty"`
	HealthCheckEndpoint string `json:"healthCheckEndpoint,omitempty"`
}

type UnixTimestamp struct {
	time.Time
}

func (ts *UnixTimestamp) UnmarshalJSON(b []byte) (err error) {
	var timestamp int64
	if err = json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	ts.Time = time.Unix(timestamp, 0)
	return
}

type QueryParam struct {
	Query          string        `json:"query,omitempty"`
	StartTimestamp UnixTimestamp `json:"start,omitempty"`
	EndTimestamp   UnixTimestamp `json:"end,omitempty"`
	Step           int64         `json:"step,omitempty"`
	LabelSelector  []string      `json:"labelSelector,omitempty"`
	Keys           []string      `json:"keys,omitempty"`
}

type QueryRequest struct {
	Source string     `json:"source"`
	Param  QueryParam `json:"params"`
}

type ResultPoint struct {
	Series string   `json:"series"`
	ValueY *float64 `json:"y"`
	ValueX *float64 `json:"x"`
}

type QueryResult struct {
	Result []ResultPoint `json:"result"`
}

const (
	DATA_SOURCE_PROM  = "prometheus"
	DATA_SOURCE_REDIS = "redis"
)
