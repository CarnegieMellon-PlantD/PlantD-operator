package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DigitalTwinSpec defines the desired state of DigitalTwin
type DigitalTwinSpec struct {
	// ModelType defines the type of the DigitalTwin model.
	ModelType string `json:"modelType,omitempty"`
	// Experiments contains the list of Experiment object references for the DigitalTwin.
	Experiments []*corev1.ObjectReference `json:"experiments,omitempty"`
}

// DigitalTwinStatus defines the observed state of DigitalTwin
type DigitalTwinStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DigitalTwin is the Schema for the digitaltwins API
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
