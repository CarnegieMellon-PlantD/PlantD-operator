package utils

import (
	"fmt"
	"net/http"
	"time"
)

const (
	timeout time.Duration = 30 * time.Second
)

func CheckHealth(url string) error {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}

	return nil
}
