package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	healthCheckTimeout time.Duration = 5 * time.Second
)

func GetHostname(serviceURL string) (string, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %v\n", err)
	}
	return u.Hostname(), nil
}

func GetPort(serviceURL string) (int32, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return -1, fmt.Errorf("Error parsing URL: %v\n", err)
	}
	if sPort := u.Port(); sPort != "" {
		port, err := strconv.Atoi(sPort)
		if err != nil {
			return -1, err
		}
		return int32(port), nil
	}
	switch u.Scheme {
	case "http":
		return 80, nil
	case "https":
		return 443, nil
	default:
		return -1, fmt.Errorf("Cannot get the default port of scheme %s", u.Scheme)
	}
}

func CheckHTTPHealth(url string) (bool, error) {
	client := http.Client{
		Timeout: healthCheckTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("expected status OK but got %d from %s", resp.StatusCode, url)
	}

	return true, nil
}

func GetHTTPPath(metricsURL string) (string, error) {
	u, err := url.Parse(metricsURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %v\n", err)
	}
	return u.Path, nil
}
