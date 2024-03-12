package v1alpha1

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PipelineInitializing string = "Initializing" // initializing includes setting up of prometheus service monitors, service setup
	PipelineAvailable    string = "Available"    // initialization has successfully completed
	PipelineEngaged      string = "Engaged"      // pipeline has been added to an experiment
	PipelineOK           string = "OK"           // if health check passed
	PipelineFailed       string = "Failed"       // if health check fails
	ProtocolHTTP         string = "http"
	ProtocolWebSocket    string = "websocket"
	ProtocolGRPC         string = "grpc"
	WithDataSet          string = "withDataSet"
	WithData             string = "withData"
)

// HTTP defines the configurations of HTTP protocol.
type HTTP struct {
	// URL defines the absolute path for an entry point of the Pipeline.
	URL string `json:"url,omitempty"`
	// Method defines the HTTP method used for the endpoint.
	Method string `json:"method,omitempty"`
	// Headers defines a map of HTTP headers.
	Headers map[string]string `json:"headers,omitempty"`
}

// WebSocket defines the configurations of websocket protocol.
type WebSocket struct {
	// Placeholder.
	URL string `json:"url,omitempty"`
	// Placeholder.
	Params map[string]string `json:"params,omitempty"`
	// Placeholder.
	Callback string `json:"callback,omitempty"`
}

// GRPC defines the configurations of gRPC protocol.
type GRPC struct {
	// Placeholder.
	Address string `json:"address,omitempty"`
	// Placeholder.
	ProtoFiles []string `json:"protoFiles,omitempty"`
	// Placeholder.
	URL string `json:"url,omitempty"`
	// Placeholder.
	Params map[string]string `json:"params,omitempty"`
	// Placeholder.
	Request map[string]string `json:"request,omitempty"`
}

// Endpoint defines the configuration of the endpoint.
type Endpoint struct {
	// Name defines the name of the endpoint. It's required when it's for pipeline endpoint.
	Name string `json:"name,omitempty"`
	// HTTP defines the configuration of the HTTP request. It's mutually exclusive with WebSocket and GRPC.
	HTTP HTTP `json:"http,omitempty"`
	// WebSocket defines the configuration of the WebSocket connection. It's mutually exclusive with HTTP and GRPC.
	WebSocket WebSocket `json:"websocket,omitempty"`
	// GRPC defines the configuration of the gRPC request. It's mutually exclusive with HTTP and WebSocket.
	GRPC GRPC `json:"grpc,omitempty"`
	// ServiceRef defines the Kubernetes Service that exposes metrics.
	ServiceRef corev1.ObjectReference `json:"serviceRef,omitempty"`
	// Internal endpoint for scraping metrics.
	monitoringv1.Endpoint `json:",inline"`
}

// SystemMetrics defines the configurations of getting system metrics.
type SystemMetrics struct {
	// Tags defines the tags for the resources of the pipeline-under-test in the cloud service provider.
	Tags map[string]string `json:"tags,omitempty"`
	// SecretRef defines the reference to the Kubernetes Secret object for authentication on the cloud service provider.
	SecretRef corev1.ObjectReference `json:"secretRef,omitempty"`
}

// MessageQueueMetrics defines the configurations of getting message queue related metrics.
type MessageQueueMetrics struct{}

// ExtraMetrics defines the configurations of getting extra metrics.
type ExtraMetrics struct {
	// System defines the configurfation of getting system metrics.
	System SystemMetrics `json:"system,omitempty"`
	// MessageQueue defines the configurfation of getting message queue related metrics.
	MessageQueue MessageQueueMetrics `json:"messageQueue,omitempty"`
}

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Endpoints for pipeline-under-test.
	PipelineEndpoints []Endpoint `json:"pipelineEndpoints"`
	// Endpoints for health check.
	HealthCheckEndpoints []string `json:"healthCheckEndpoints,omitempty"`
	// Endpoints for metrics.
	MetricsEndpoint Endpoint `json:"metricsEndpoint,omitempty"`
	// Extra metrics, such as CPU utilzation, I/O and etc.
	ExtraMetrics ExtraMetrics `json:"extraMetrics,omitempty"`
	// In cluster flag. True indecates the pipeline-under-test is deployed in the same cluster as the plantD. Otherwise it should be False.
	InCluster bool `json:"inCluster,omitempty"`
	// State which cloud service provider the pipeline is deployed.
	CloudVendor string `json:"cloudVendor,omitempty"`
	// Cost calculation flag.
	EnableCostCalculation bool `json:"enableCostCalculation,omitempty"`
	// Internal usage. For experiment object to lock the pipeline object.
	ExperimentRef corev1.ObjectReference `json:"experimentRef,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// PipelineState defines the state of the Pipeline.
	PipelineState string `json:"pipelineState,omitempty"`
	// StatusCheck defines the health status of the Pipeline.
	StatusCheck string `json:"statusCheck,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="PipelineState",type="string",JSONPath=".status.pipelineState"
//+kubebuilder:printcolumn:name="StatusCheck",type="string",JSONPath=".status.statusCheck"

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the Pipeline.
	Spec PipelineSpec `json:"spec,omitempty"`
	// Status defines the status of the Pipeline.
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of Pipelines.
	Items []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
