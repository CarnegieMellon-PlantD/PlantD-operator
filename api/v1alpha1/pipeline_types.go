package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PipelineAvailability defines the availability of the Pipeline.
type PipelineAvailability string

const (
	PipelineReady PipelineAvailability = "Ready"
	PipelineInUse PipelineAvailability = "In-Use"
)

// HTTP defines the configurations of HTTP protocol in endpoint.
type HTTP struct {
	// URL of the HTTP request.
	URL string `json:"url"`
	// Method of the HTTP request.
	Method string `json:"method"`
	// Headers of the HTTP request.
	Headers map[string]string `json:"headers,omitempty"`
}

// PipelineEndpoint defines the endpoint for data ingestion in Pipeline.
type PipelineEndpoint struct {
	// Name of the endpoint.
	Name string `json:"name"`
	// Configurations of the HTTP protocol.
	HTTP *HTTP `json:"http"`
}

// MetricsEndpoint defines the endpoint for metrics scraping in Pipeline.
type MetricsEndpoint struct {
	// Configurations of the HTTP protocol.
	// Only the scheme, host, and port of the `http.url` field will be used.
	// Must be set if `inCluster` is set to `false` in the Pipeline.
	HTTP *HTTP `json:"http,omitempty"`
	// Reference to the Service.
	// The Service must be in the same namespace as the Pipeline.
	// Must be set if `inCluster` is set to `true` in the Pipeline.
	ServiceRef *corev1.LocalObjectReference `json:"serviceRef,omitempty"`
	// Name of the Service port to use.
	// Default to "metrics".
	Port string `json:"port,omitempty"`
	// Path of the endpoint.
	// Default to "/metrics".
	Path string `json:"path,omitempty"`
}

// PipelineSpec defines the desired state of Pipeline.
type PipelineSpec struct {
	// Whether the Pipeline is deployed within the cluster or not.
	// When set to `false`, Services of type ExternalName will be created to access the Pipeline.
	// When set to `true`, the Pipeline will be accessed by its Services.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	InCluster bool `json:"inCluster,omitempty"`
	// List of endpoints for data ingestion.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	PipelineEndpoints []PipelineEndpoint `json:"pipelineEndpoints"`
	// Endpoint for metrics scraping.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	MetricsEndpoint *MetricsEndpoint `json:"metricsEndpoint,omitempty"`
	// List of URLs for health check.
	// An HTTP GET request will be made to each URL, and all of them should return 200 OK to pass the health check.
	// If the list is empty, no health check will be performed.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	HealthCheckURLs []string `json:"healthCheckURLs,omitempty"`
	// Whether to enable cost calculation for the Pipeline.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	EnableCostCalculation bool `json:"enableCostCalculation,omitempty"`
	// Cloud provider of the Pipeline. Available values are `aws`, `azure`, and `gcp`.
	// +kubebuilder:validation:Enum=aws;azure;gcp
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	CloudProvider string `json:"cloudProvider,omitempty"`
	// Map of tags to select cloud resources of the Pipeline. Equivalent to the tags in the cloud service provider.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="value is immutable"
	Tags map[string]string `json:"tags,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline.
type PipelineStatus struct {
	// Availability of the Pipeline.
	Availability PipelineAvailability `json:"availability,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Availability",type="string",JSONPath=".status.availability"
//+kubebuilder:printcolumn:name="Liveness",type="string",JSONPath=".status.liveness"

// Pipeline is the Schema for the pipelines API
// +kubebuilder:validation:XValidation:rule="size(self.metadata.name) <= 55",message="must contain at most 55 characters"
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
