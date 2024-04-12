package utils

import (
	"fmt"
	"net/url"
	"strconv"
)

// GetURLHostname returns the hostname of the URL.
func GetURLHostname(in string) (string, error) {
	u, err := url.Parse(in)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// GetURLPort returns the port of the URL.
func GetURLPort(in string) (int, error) {
	u, err := url.Parse(in)
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
		return 0, fmt.Errorf("unsupported scheme \"%s\" in URL", u.Scheme)
	}
}
