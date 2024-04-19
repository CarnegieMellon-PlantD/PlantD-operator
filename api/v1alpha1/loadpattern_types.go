package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Stage defines how the load ramps up or down.
type Stage struct {
	// Target load to reach at the end of the stage.
	// Equivalent to the "ramping-arrival-rate" executor's `stages[].target` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	// +kubebuilder:validation:Minimum=0
	Target int64 `json:"target"`
	// Duration of the stage, also the time to reach the target load.
	// Equivalent to the "ramping-arrival-rate" executor's `stages[].duration` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	Duration string `json:"duration"`
}

// LoadPatternSpec defines the desired state of LoadPattern.
type LoadPatternSpec struct {
	// List of stages in the LoadPattern.
	// Equivalent to the "ramping-arrival-rate" executor's `stages` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	// +kubebuilder:validation:MinItems=1
	Stages []Stage `json:"stages"`
	// Number of VUs to pre-allocate before Experiment start.
	// Equivalent to the "ramping-arrival-rate" executor's `preAllocatedVUs` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	// +kubebuilder:validation:Minimum=0
	PreAllocatedVUs int64 `json:"preAllocatedVUs,omitempty"`
	// Number of requests per `timeUnit` period at Experiment start.
	// Equivalent to the "ramping-arrival-rate" executor's `startRate` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	// +kubebuilder:validation:Minimum=0
	StartRate int64 `json:"startRate"`
	// Period of time to apply to the `startRate` and `stages[].target` fields.
	// Equivalent to the "ramping-arrival-rate" executor's `timeUnit` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	TimeUnit string `json:"timeUnit,omitempty"`
	// Maximum number of VUs to allow for allocation during Experiment.
	// Equivalent to the "ramping-arrival-rate" executor's `maxVUs` option in K6.
	// See https://k6.io/docs/using-k6/scenarios/executors/ramping-arrival-rate/#options for more details.
	// +kubebuilder:validation:Minimum=0
	MaxVUs int64 `json:"maxVUs,omitempty"`
}

// LoadPatternStatus defines the observed state of LoadPattern.
type LoadPatternStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LoadPattern is the Schema for the loadpatterns API
type LoadPattern struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadPatternSpec   `json:"spec,omitempty"`
	Status LoadPatternStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LoadPatternList contains a list of LoadPattern
type LoadPatternList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadPattern `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoadPattern{}, &LoadPatternList{})
}
