package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimulationSpec defines the desired state of Simulation
type SimulationSpec struct {
	// DigitalTwin object for the Simulation.
	DigitalTwinRef *corev1.ObjectReference `json:"digitalTwinRef"`
	// TrafficModel object for the Simulation.
	TrafficModelRef *corev1.ObjectReference `json:"trafficModelRef"`
	// NetCost object for the Simulation.
	// Optional.
	NetCostRef *corev1.ObjectReference `json:"netCostRef,omitempty"`
	// Scenario object for the Simulation.
	// The task names in the Scenario must be the name of a Schema in the DataSet used by the DigitalTwin.
	// Mandatory if the `digitalTwinType` field of the DigitalTwin is `schemaaware`.
	// Always ignored otherwise.
	ScenarioRef *corev1.ObjectReference `json:"scenarioRef,omitempty"`
}

// SimulationStatus defines the observed state of Simulation
type SimulationStatus struct{}

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
