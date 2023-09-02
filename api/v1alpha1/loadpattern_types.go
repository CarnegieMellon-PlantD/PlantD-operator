package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Stage defines the stage configuration of the load.
type Stage struct {
	// Target defines the target requests per second.
	Target int `json:"target"`
	// Duration defines the duration of the current stage.
	Duration string `json:"duration"`
}

// LoadPatternSpec defines the desired state of LoadPattern
type LoadPatternSpec struct {
	// Stages defines a list of stages for the LoadPattern.
	Stages []Stage `json:"stages,omitempty"`
	// PreAllocatedVUs defines pre-allocated virtual users for the K6 load generator.
	PreAllocatedVUs int `json:"preAllocatedVUs,omitempty"`
	// StartRate defines the initial requests per second when the K6 load generator starts.
	StartRate int `json:"startRate,omitempty"`
	// MaxVUs defines the maximum virtual users for the K6 load generator.
	MaxVUs int `json:"maxVUs,omitempty"`
	// TimeUnit defines the unit of the time for K6 load generator.
	TimeUnit string `json:"timeUnit,omitempty"`
}

// LoadPatternStatus defines the observed state of LoadPattern
type LoadPatternStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LoadPattern is the Schema for the loadpatterns API
type LoadPattern struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specification of the LoadPattern.
	Spec LoadPatternSpec `json:"spec,omitempty"`
	// Status defines the status of the LoadPattern.
	Status LoadPatternStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LoadPatternList contains a list of LoadPattern
type LoadPatternList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items defines a list of LoadPattern.
	Items []LoadPattern `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadPattern{}, &LoadPatternList{})
}
