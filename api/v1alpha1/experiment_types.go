package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LoadPatternConfig defines the configuration of the load pattern in the experiment.
type LoadPatternConfig struct {
	// EndpointName defines the name of endpoint where to send the requests.
	// It should match the name of endpoint declared in the specification of the pipeline.
	EndpointName string `json:"endpointName"`
	// LoadPatternRef defines s reference of the LoadPattern object.
	LoadPatternRef corev1.ObjectReference `json:"loadPatternRef"`
}

// ExperimentSpec defines the desired state of Experiment
type ExperimentSpec struct {
	// PipelineRef defines s reference of the Pipeline object.
	PipelineRef corev1.ObjectReference `json:"pipelineRef"`
	// LoadPatterns defines a list of configuration of name of endpoints and LoadPatterns.
	LoadPatterns []LoadPatternConfig `json:"loadPatterns"`
	// ScheduledTime defines the scheduled time for the Experiment.
	ScheduledTime metav1.Time `json:"scheduledTime,omitempty"`
}

// ExperimentStatus defines the observed state of Experiment
type ExperimentStatus struct {
	// ExperimentState defines the state of the Experiment.
	ExperimentState string `json:"experimentState,omitempty"`
	// Protocols defines a map of name of endpoint (key) to request protocol (value).
	Protocols map[string]string `json:"protocols,omitempty"`
	// Tags defines the a map of key-value pair that use for tagging cloud resources.
	Tags map[string]string `json:"tags,omitempty"`
	// Duration defines the duration of the K6 load generator.
	Duration map[string]metav1.Duration `json:"duration,omitempty"`
	// StartTime defines the start of the K6 load generator.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// EndTime defines the end of the Experiment.
	// TODO: Add microservice to calculate the end time of the experiment.
	EndTime *metav1.Time `json:"endTime,omitempty"`
	// CloudVendor defines the cloud service provider which the pipeline-under-test is deployed.
	CloudVendor string `json:"cloudVendor,omitempty"`
	// EnableCostCalculation defines teh flag of cost calculation.
	EnableCostCalculation bool `json:"enableCostCalculation,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Experiment is the Schema for the experiments API
type Experiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specification of the Experiment.
	Spec ExperimentSpec `json:"spec,omitempty"`
	// Status defines the status of the Experiment.
	Status ExperimentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExperimentList contains a list of Experiment
type ExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items defines a list of Experiment.
	Items []Experiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experiment{}, &ExperimentList{})
}
