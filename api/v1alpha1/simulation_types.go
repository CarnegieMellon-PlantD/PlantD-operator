package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SimulationJobStatus defines the status of the digital twin job.
type SimulationJobStatus string

const (
	SimulationRunning   SimulationJobStatus = "Running"
	SimulationCompleted SimulationJobStatus = "Completed"
	SimulationFailed    SimulationJobStatus = "Failed"
)

// SimulationSpec defines the desired state of Simulation
type SimulationSpec struct {
	// Container image to use for the simulation.
	Image string `json:"image,omitempty"`
	// DigitalTwin object for the Simulation.
	DigitalTwinRef *corev1.ObjectReference `json:"digitalTwinRef,omitempty"`
	// TrafficModel object for the Simulation.
	TrafficModelRef *corev1.ObjectReference `json:"trafficModelRef"`
	// NetCost object for the Simulation.
	// Optional if the `digitalTwinType` field is specified and the target DigitalTwin is of type `schemaaware`.
	// Always ignored otherwise.
	NetCostRef *corev1.ObjectReference `json:"netCostRef,omitempty"`
	// Scenario object for the Simulation.
	// The task names in the Scenario must be the name of a Schema in the DataSet used by the DigitalTwin.
	// Required if the `digitalTwinType` field is specified and the target DigitalTwin is of type `schemaaware`.
	// Always ignored otherwise.
	ScenarioRef *corev1.ObjectReference `json:"scenarioRef,omitempty"`
}

// SimulationStatus defines the observed state of Simulation
type SimulationStatus struct {
	// Status of the digital twin job.
	JobStatus SimulationJobStatus `json:"jobStatus,omitempty"`
	// Error message.
	Error string `json:"error,omitempty"`
}

// The name of the Pod of the Simulation will be
// "<simulation-name>-sim-<random 5 chars>".
// So, we have 53 characters for the name to meet the 63-character limit.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="JobStatus",type="string",JSONPath=".status.jobStatus"

// Simulation is the Schema for the simulations API
// +kubebuilder:validation:XValidation:rule="size(self.metadata.name) <= 53",message="must contain at most 53 characters"
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
