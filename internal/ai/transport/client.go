package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http    *http.Client
	retries int
}

type Config struct {
	Timeout    time.Duration
	MaxRetries int
	Proxy      string
}

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

func New(cfg Config) (*Client, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 3
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,

		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}

	if cfg.Proxy != "" {
		u, err := url.Parse(cfg.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy: %w", err)
		}

		transport.Proxy = http.ProxyURL(u)
	}

	return &Client{
		http: &http.Client{
			Transport: transport,
			Timeout:   cfg.Timeout,
		},
		retries: cfg.MaxRetries,
	}, nil
}

// DoWithContext executes an HTTP request with context (نسخه با context)
func (c *Client) DoWithContext(ctx context.Context, req *Request) ([]byte, error) {
	var lastErr error
	maxRetries := c.retries
	if maxRetries < 0 {
		maxRetries = 3
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var bodyReader io.Reader
		if req.Body != nil {
			bodyReader = bytes.NewReader(req.Body)
		}

		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}

		resp, err := c.http.Do(httpReq)
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				backoff := time.Duration(1<<uint(attempt)) * time.Second
				time.Sleep(backoff)
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				backoff := time.Duration(1<<uint(attempt)) * time.Second
				time.Sleep(backoff)
				continue
			}
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return body, nil
		}

		if resp.StatusCode == 429 {
			lastErr = fmt.Errorf("rate limited (status %d): %s", resp.StatusCode, string(body))
			if attempt < maxRetries {
				backoff := time.Duration(1<<uint(attempt+1)) * time.Second
				time.Sleep(backoff)
				continue
			}
			return nil, lastErr
		}

		lastErr = fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		if attempt < maxRetries {
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			time.Sleep(backoff)
			continue
		}
		return nil, lastErr
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, lastErr)
}
