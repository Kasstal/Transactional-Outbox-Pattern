package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Config stores configuration for the API client
type ClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

// Client is the API client structure
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new instance of the API client
func NewClient(config ClientConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL: config.BaseURL,
	}
}

// SendRequest sends an HTTP request to another service in the server
func (c *Client) SendRequest(endpoint string, method string, body interface{}) (*http.Response, error) {
	// Prepare URL
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	// Prepare the body
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Handle non 200 HTTP status codes
	if resp.StatusCode >= 400 {
		return resp, errors.New(fmt.Sprintf("received non-2xx response: %d", resp.StatusCode))
	}

	return resp, nil
}
