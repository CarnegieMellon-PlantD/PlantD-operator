package utils

// IntRange defines a range using two integers as boundaries.
type IntRange struct {
	// Minimum value of the range.
	Min int `json:"min"`
	// Maximum value of the range.
	Max int `json:"max"`
}
