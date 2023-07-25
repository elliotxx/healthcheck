package checks

import (
	"io"
	"net/http"
	"time"
)

// NewCustomPingCheck creates a new custom health check that pings a
// specified URL.
func NewCustomPingCheck(url, method string, timeout int, body io.Reader, headers map[string]string) Check {
	// Default to GET method and 500 millisecond timeout if not
	// specified.
	if method == "" {
		method = http.MethodGet
	}

	if timeout == 0 {
		timeout = 500
	}

	customPingCheck := &customPingCheck{
		url:     url,
		method:  method,
		timeout: timeout,
		body:    body,
		headers: headers,
	}

	// Create an HTTP client with the specified timeout.
	customPingCheck.client = http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	return customPingCheck
}

// customPingCheck represents a health check that pings a specified
// URL.
type customPingCheck struct {
	url     string
	method  string
	timeout int
	client  http.Client
	body    io.Reader
	headers map[string]string
}

// Pass checks if the specified URL is reachable.
func (p *customPingCheck) Pass() bool {
	req, err := http.NewRequest(p.method, p.url, p.body)
	if err != nil {
		return false
	}

	for key, value := range p.headers {
		req.Header.Add(key, value)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < http.StatusMultipleChoices
}

// Name returns the name of the custom ping check.
func (p *customPingCheck) Name() string {
	return "Ping-" + p.url
}
