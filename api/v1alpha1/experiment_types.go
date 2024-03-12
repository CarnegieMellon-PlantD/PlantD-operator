package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataSpec defines the data to be sent to the endpoint.
type DataSpec struct {
	// PlainText defines a plain text data.
	PlainText string `json:"plainText,omitempty"`
	// DataSetRef defines the reference of the DataSet object.
	DataSetRef corev1.ObjectReference `json:"dataSetRef,omitempty"`
}

// EndpointSpec defines the DataSet and LoadPattern to be used for an endpoint.
type EndpointSpec struct {
	// EndpointName defines the name of endpoint.
	// It should be the name of an existing endpoint defined in the Pipeline used in the Experiment.
	EndpointName string `json:"endpointName,omitempty"`
	// DataSpec defines the data to be sent to the endpoint.
	DataSpec DataSpec `json:"dataSpec,omitempty"`
	// LoadPatternRef defines the reference of the LoadPattern object.
	LoadPatternRef corev1.ObjectReference `json:"loadPatternRef,omitempty"`
}

// ExperimentSpec defines the desired state of Experiment
type ExperimentSpec struct {
	// PipelineRef defines a reference of the Pipeline object.
	PipelineRef corev1.ObjectReference `json:"pipelineRef,omitempty"`
	// EndpointSpecs defines a list of configurations for the endpoints.
	EndpointSpecs []EndpointSpec `json:"endpointSpecs,omitempty"`
	// ScheduledTime defines the scheduled time for the Experiment.
	ScheduledTime metav1.Time `json:"scheduledTime,omitempty"`
}

// ExperimentStatus defines the observed state of Experiment
type ExperimentStatus struct {
	// ExperimentState defines the state of the Experiment.
	ExperimentState string `json:"experimentState,omitempty"`
	// Protocols defines the map of name of endpoint (key) to request protocol (value).
	Protocols map[string]string `json:"protocols,omitempty"`
	// Tags defines the map of key-value pair that use for tagging cloud resources.
	Tags map[string]string `json:"tags,omitempty"`
	// Duration defines the duration of the K6 load generator.
	Duration map[string]metav1.Duration `json:"duration,omitempty"`
	// StartTime defines the start of the K6 load generator.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// EndTime defines the end of the Experiment.
	EndTime *metav1.Time `json:"endTime,omitempty"`
	// CloudVendor defines the cloud service provider which the pipeline-under-test is deployed.
	CloudVendor string `json:"cloudVendor,omitempty"`
	// EnableCostCalculation defines the flag of cost calculation.
	EnableCostCalculation bool `json:"enableCostCalculation,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="ExperimentState",type="string",JSONPath=".status.experimentState"
//+kubebuilder:printcolumn:name="Duration",type="string",JSONPath=".status.duration"
//+kubebuilder:printcolumn:name="StartTime",type="string",JSONPath=".status.startTime"

// Experiment is the Schema for the experiments API
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the Experiment.
	Spec ExperimentSpec `json:"spec,omitempty"`
	// Status defines the status of the Experiment.
	Status ExperimentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExperimentList contains a list of Experiments.
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of Experiments.
	Items []Experiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}
