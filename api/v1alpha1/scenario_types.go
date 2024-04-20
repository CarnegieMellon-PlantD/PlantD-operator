package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScenarioTask defines the task in the Scenario.
type ScenarioTask struct {
	// Name of the task. Should be a Schema name.
	Name string `json:"name"`
	// The size of a single upload in bytes.
	Size *resource.Quantity `json:"size"`
	// Range of the number range of the devices to send the data.
	SendingDevices NaturalIntRange `json:"sendingDevices"`
	// Range of the frequency of data pushes per month.
	PushFrequencyPerMonth NaturalIntRange `json:"pushFrequencyPerMonth"`
	// List of months the task will apply to.
	// For example, `[1, 12]` means the task will apply to January and December.
	MonthsRelevant []int `json:"monthsRelevant"`
}

// ScenarioSpec defines the desired state of Scenario
type ScenarioSpec struct {
	// List of tasks in the Scenario.
	// +kubebuilder:validation:MinItems=1
	Tasks []ScenarioTask `json:"tasks"`
}

// ScenarioStatus defines the observed state of Scenario
type ScenarioStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scenario is the Schema for the scenarios API
type Scenario struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScenarioSpec   `json:"spec,omitempty"`
	Status ScenarioStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScenarioList contains a list of Scenario
type ScenarioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scenario `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scenario{}, &ScenarioList{})
}
