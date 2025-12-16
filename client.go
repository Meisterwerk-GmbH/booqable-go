package bq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the public Booqable API endpoint.
	DefaultBaseURL = "https://api.booqable.com/v4"
	defaultTimeout = 30 * time.Second
	defaultAgent   = "bq-go/0.1"
)

// Client talks to the Booqable API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	UserAgent  string
}

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		if baseURL != "" {
			c.BaseURL = baseURL
		}
	}
}

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		if hc != nil {
			c.HTTPClient = hc
		}
	}
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(agent string) ClientOption {
	return func(c *Client) {
		if agent != "" {
			c.UserAgent = agent
		}
	}
}

// NewClient creates a Client with sane defaults.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:    DefaultBaseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		UserAgent:  defaultAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get issues a GET request and decodes JSON into v.
func (c *Client) Get(ctx context.Context, p string, query url.Values, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, p, query, nil)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

func (c *Client) newRequest(ctx context.Context, method, p string, query url.Values, body io.Reader) (*http.Request, error) {
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	if query == nil {
		query = url.Values{}
	}
	base.Path = path.Join(strings.TrimSuffix(base.Path, "/"), strings.TrimPrefix(p, "/"))
	base.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, base.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.APIKey != "" {
		req.Header.Set("X-Api-Key", c.APIKey)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	hc := c.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: defaultTimeout}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}
	if v == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// APIError represents a non-2xx response.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: status=%d body=%s", e.StatusCode, e.Body)
}
