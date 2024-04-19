package v1alpha1

// NaturalIntRange defines a range using two non-negative integers as boundaries.
type NaturalIntRange struct {
	// Minimum value of the range.
	// +kubebuilder:validation:Minimum=0
	Min int32 `json:"min"`
	// Maximum value of the range.
	// +kubebuilder:validation:Minimum=0
	Max int32 `json:"max"`
}
