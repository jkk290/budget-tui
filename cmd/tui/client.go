package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	jwt        string
}

func newClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SetJWT(token string) {
	c.jwt = token
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+c.jwt)
	}

	return req, nil
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, v any) (*http.Request, error) {
	var body io.Reader

	if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}
