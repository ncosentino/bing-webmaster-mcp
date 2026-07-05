// Package indexnow provides a client for the Bing IndexNow protocol.
package indexnow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const httpTimeout = 30 * time.Second

// apiBaseURL is the IndexNow endpoint.
// It is a variable so tests can override it to point at a local test server.
var apiBaseURL = "https://www.bing.com/indexnow"

// Client calls the IndexNow endpoint.
type Client struct {
	httpClient *http.Client
	defaultKey string
}

// NewClient creates a Client with an optional default IndexNow key.
func NewClient(defaultKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: httpTimeout},
		defaultKey: defaultKey,
	}
}

type submitResult struct {
	Host        string    `json:"host"`
	URLList     []string  `json:"urlList"`
	KeyLocation string    `json:"keyLocation,omitempty"`
	Success     bool      `json:"success"`
	KeySource   string    `json:"keySource"`
	SubmittedAt time.Time `json:"submittedAt"`
}

type submitRequest struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation,omitempty"`
	URLList     []string `json:"urlList"`
}

// SubmitURLs submits one or more URLs to IndexNow.
func (c *Client) SubmitURLs(ctx context.Context, host string, urlList []string, keyOverride string, keyLocation string) (*submitResult, error) {
	if len(urlList) == 0 {
		return nil, fmt.Errorf("urlList must contain at least one URL")
	}

	keySource := "override"
	key := keyOverride
	if key == "" {
		key = c.defaultKey
		keySource = "configured"
	}
	if key == "" {
		return nil, fmt.Errorf("no IndexNow key configured; provide the key parameter or configure BING_INDEXNOW_KEY / --indexnow-key")
	}

	payload := submitRequest{
		Host:        host,
		Key:         key,
		KeyLocation: keyLocation,
		URLList:     urlList,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("IndexNow returned HTTP %d: %s", resp.StatusCode, statusReason(resp.StatusCode, string(responseBody)))
	}

	return &submitResult{
		Host:        host,
		URLList:     urlList,
		KeyLocation: keyLocation,
		Success:     true,
		KeySource:   keySource,
		SubmittedAt: time.Now().UTC(),
	}, nil
}

func statusReason(statusCode int, body string) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "bad request"
	case http.StatusForbidden:
		return "key invalid"
	case http.StatusUnprocessableEntity:
		return "URLs do not belong to host or key location does not match"
	case http.StatusTooManyRequests:
		return "too many requests"
	default:
		if body != "" {
			return truncate(body, 300)
		}
		return http.StatusText(statusCode)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
