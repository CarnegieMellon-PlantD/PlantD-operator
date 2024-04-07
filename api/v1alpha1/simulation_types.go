package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimulationSpec defines the desired state of Simulation
type SimulationSpec struct {
	// TrafficModelRef defines the TrafficModel object reference for the Simulation.
	TrafficModelRef corev1.ObjectReference `json:"trafficModelRef"`
	// DigitalTwinRef defines the DigitalTwin object reference for the Simulation.
	DigitalTwinRef corev1.ObjectReference `json:"digitalTwinRef"`
}

// SimulationStatus defines the observed state of Simulation
type SimulationStatus struct {
	// DigitalTwinState is the state of the DigitalTwin.
	DigitalTwinState string `json:"digitalTwinState,omitempty"`
	// TrafficModelState is the state of the TrafficModel.
	TrafficModelState string `json:"trafficModalState,omitempty"`
	// PodName is the pod name of the digital twin job.
	PodName string `json:"podName,omitempty"`
	// JobStatus is the status of the digital twin job.
	JobStatus string `json:"jobStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="DigitalTwinState",type="string",JSONPath=".status.digitaltwinState"
//+kubebuilder:printcolumn:name="TrafficModelState",type="string",JSONPath=".status.trafficmodalState"
//+kubebuilder:printcolumn:name="PodName",type="string",JSONPath=".status.podName"
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"

// Simulation is the Schema for the simulations API
type Simulation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimulationSpec   `json:"spec,omitempty"`
	Status SimulationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SimulationList contains a list of Simulation
type SimulationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Simulation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Simulation{}, &SimulationList{})
}
