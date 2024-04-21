package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DigitalTwinJobStatus defines the status of the Experiments created by DigitalTwin.
type DigitalTwinJobStatus string

const (
	DigitalTwinRunning   DigitalTwinJobStatus = "Running"
	DigitalTwinCompleted DigitalTwinJobStatus = "Completed"
	DigitalTwinFailed    DigitalTwinJobStatus = "Failed"
)

// DigitalTwinSpec defines the desired state of DigitalTwin.
type DigitalTwinSpec struct {
	// Type of digital twin model.
	// Available values are `simple`, `quickscaling`, and `autoscaling`.
	// +kubebuilder:validation:Enum=simple;quickscaling;autoscaling
	ModelType string `json:"modelType"`
	// Type of digital twin.
	// Available values are `regular` and `schemaaware`.
	// +kubebuilder:validation:Enum=regular;schemaaware
	DigitalTwinType string `json:"digitalTwinType"`
	// Existing Experiments to retrieve metrics data from to train the DigitalTwin.
	// Effective only when `digitalTwinType` is `regular`.
	Experiments []*corev1.ObjectReference `json:"experiments,omitempty"`
	// DataSet to break down into Schemas to train the DigitalTwin.
	// Effective only when `digitalTwinType` is `schemaaware`.
	DataSet *corev1.LocalObjectReference `json:"dataSet,omitempty"`
	// Pipeline to use to train the DigitalTwin.
	// Effective only when `digitalTwinType` is `schemaaware`.
	Pipeline *corev1.LocalObjectReference `json:"pipeline,omitempty"`
	// Maximum RPS in the populated LoadPatterns.
	// Effective only when `digitalTwinType` is `schemaaware`.
	PipelineCapacity int32 `json:"pipelineCapacity,omitempty"`
}

// DigitalTwinStatus defines the observed state of DigitalTwin.
type DigitalTwinStatus struct {
	// Status of the Experiments created by DigitalTwin.
	JobStatus DigitalTwinJobStatus `json:"jobStatus,omitempty"`
	// Error message.
	Error string `json:"error,omitempty"`
}

// The name of the Experiment for the DigitalTwin will be
// "<digitaltwin-name>-pure-<up to 4 digits of schema index>".
// So, we have 22 characters for the name to meet the 32-character limit.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"

// DigitalTwin is the Schema for the digitaltwins API
// +kubebuilder:validation:XValidation:rule="size(self.metadata.name) <= 22",message="must contain at most 22 characters"
type DigitalTwin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DigitalTwinSpec   `json:"spec,omitempty"`
	Status DigitalTwinStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DigitalTwinList contains a list of DigitalTwin
type DigitalTwinList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DigitalTwin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DigitalTwin{}, &DigitalTwinList{})
}
