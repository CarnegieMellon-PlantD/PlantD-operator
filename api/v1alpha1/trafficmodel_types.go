package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TrafficModelSpec defines the desired state of TrafficModel.
type TrafficModelSpec struct {
	// TrafficModel configuration in JSON.
	Config string `json:"config"`
}

// TrafficModelStatus defines the observed state of TrafficModel.
type TrafficModelStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TrafficModel is the Schema for the trafficmodels API
type TrafficModel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specifications of the TrafficModel.
	Spec TrafficModelSpec `json:"spec,omitempty"`
	// Status defines the status of the TrafficModel.
	Status TrafficModelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TrafficModelList contains a list of TrafficModel
type TrafficModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items defines a list of TrafficModels.
	Items []TrafficModel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrafficModel{}, &TrafficModelList{})
}
