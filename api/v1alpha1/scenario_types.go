package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataSetConfig defines the parameters to generate DataSet
type DataSetConfig struct {
	CompressPerSchema    bool   `json:"compressPerSchema,omitempty"`
	CompressedFileFormat string `json:"compressedFileFormat,omitempty"`
	FileFormat           string `json:"fileFormat,omitempty"`
}

// ScenarioTask defines the task to be executed in the Scenario
type ScenarioTask struct {
	// Name defines the name of the task.
	Name string `json:"name,omitempty"`
	// Size defines the size of a single upload in bytes.
	Size resource.Quantity `json:"size,omitempty"`
	// SendingDevices defines the range of the devices to send the data.
	SendingDevices map[string]int `json:"sendingDevices,omitempty"`
	// PushFrequencyPerMonth defines the range of how many times the data is pushed per month.
	PushFrequencyPerMonth map[string]int `json:"pushFrequencyPerMonth,omitempty"`
	// MonthsRelevant defines the months the task is relevant.
	MonthsRelevant []int `json:"monthsRelevant,omitempty"`
}

// ScenarioSpec defines the desired state of Scenario
type ScenarioSpec struct {
	// DataSetConfig defines the parameters to generate DataSet.
	DataSetConfig DataSetConfig `json:"dataSetConfig"`
	// PipelineRef defines the reference to the Pipeline object.
	PipelineRef corev1.ObjectReference `json:"pipelineRef"`
	// Tasks defines the list of tasks to be executed in the Scenario.
	Tasks []ScenarioTask `json:"tasks,omitempty"`
}

// ScenarioStatus defines the observed state of Scenario
type ScenarioStatus struct {
	// IsPopulated defines whether the Scenario has populated the resources.
	IsPopulated bool `json:"isPopulated,omitempty"`
}

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
