package utils

import (
	"fmt"
	"net/http"
)

func PingURL(url string) error {
	// Send an HTTP GET request to the specified URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to ping %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to ping %s: %s", url, resp.Status)
	}

	return nil
}
