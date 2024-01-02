package requester

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/bugrakocabay/konsume/pkg/metrics"
)

// HTTPRequester is an interface for sending HTTP requests, useful for mocking in tests
type HTTPRequester interface {
	SendRequest() (*http.Response, error)
}

// Requester is the struct that contains the request information.
type Requester struct {
	Endpoint string
	Method   string
	Body     []byte
	Headers  map[string]string
}

// NewRequester creates a new Requester struct.
func NewRequester(endpoint, method string, body []byte, headers map[string]string) *Requester {
	return &Requester{
		Endpoint: endpoint,
		Method:   method,
		Body:     body,
		Headers:  headers,
	}
}

// SendRequest sends the request to the given endpoint.
func (r *Requester) SendRequest() (*http.Response, error) {
	var (
		resp *http.Response
		err  error
		body io.Reader
	)

	if len(r.Body) > 0 {
		body = bytes.NewBuffer(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.Endpoint, body)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, err
	}
	if len(r.Headers) > 0 {
		for k, v := range r.Headers {
			req.Header.Add(k, v)
		}
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Failed to send request", "error", err)
		return resp, err
	}
	defer resp.Body.Close()

	metrics.HttpRequestsMade.Inc()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		metrics.HttpRequestsSucceeded.Inc()
	} else {
		metrics.HttpRequestsFailed.Inc()
	}
	return resp, nil
}
