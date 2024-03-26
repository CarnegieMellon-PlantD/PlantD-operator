package utils

// IntRange defines a range using two non-negative integers as boundaries.
type IntRange struct {
	// Minimum value of the range.
	// +kubebuilder:validation:Minimum=0
	Min int `json:"min"`
	// Maximum value of the range.
	// +kubebuilder:validation:Minimum=0
	Max int `json:"max"`
}
