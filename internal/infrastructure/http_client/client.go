package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Client is a wrapper around http.Client with additional functionality
type Client struct {
	client *http.Client
}

// RequestConfig holds configuration for an HTTP request
type RequestConfig struct {
	Headers map[string]string
	Timeout time.Duration
}

// NewClient creates a new HTTP client with default timeout
func NewClient(defaultTimeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Get performs an HTTP GET request
func (c *Client) Get(ctx context.Context, url string, config *RequestConfig) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req, config)
}

// Post performs an HTTP POST request with a JSON body
func (c *Client) Post(ctx context.Context, url string, body interface{}, config *RequestConfig) (*http.Response, error) {
	// Marshal body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, jsonReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, config)
}

// do performs the HTTP request with the given configuration
func (c *Client) do(req *http.Request, config *RequestConfig) (*http.Response, error) {
	// Use client's copy to avoid modifying the original
	client := *c.client

	// Apply custom timeout if provided
	if config != nil && config.Timeout > 0 {
		client.Timeout = config.Timeout
	}

	// Apply headers if provided
	if config != nil && len(config.Headers) > 0 {
		for key, value := range config.Headers {
			req.Header.Set(key, value)
		}
	}

	// Execute request
	return client.Do(req)
}

// UnmarshalResponse reads the response body and unmarshals it into the given target
func UnmarshalResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

// Helper function to convert JSON bytes to an io.Reader
func jsonReader(jsonBytes []byte) io.Reader {
	return io.NopCloser(bytes.NewReader(jsonBytes))
}