package v1alpha1

import (
	"fmt"
	"net/url"
	"strconv"
)

// IntRange defines a range using two non-negative integers as boundaries.
type IntRange struct {
	// Minimum value of the range.
	// +kubebuilder:validation:Minimum=0
	Min int `json:"min"`
	// Maximum value of the range.
	// +kubebuilder:validation:Minimum=0
	Max int `json:"max"`
}

// GetHostname returns the hostname of the URL.
func (h *HTTP) GetHostname() (string, error) {
	u, err := url.Parse(h.URL)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// GetPort returns the port of the URL.
func (h *HTTP) GetPort() (int, error) {
	u, err := url.Parse(h.URL)
	if err != nil {
		return 0, err
	}
	if sPort := u.Port(); sPort != "" {
		port, err := strconv.Atoi(sPort)
		if err != nil {
			return 0, err
		}
		return port, nil
	}

	// If no port is explicitly specified, return the default port based on the scheme
	switch u.Scheme {
	case "http":
		return 80, nil
	case "https":
		return 443, nil
	default:
		return 0, fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
}
