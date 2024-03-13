package e2e

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"

	"gopkg.in/yaml.v3"
)

type RequestCapture struct {
	Mutex            sync.Mutex
	ReceivedRequests []HTTPRequestExpectation
}

// writeConfigToFile writes the given config to a temporary file and returns the file path and a cleanup function
func writeConfigToFile(cfg *config.Config) (string, func()) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "konsume-config-*.yaml")
	if err != nil {
		log.Fatalf("Failed to create temp file for config: %v", err)
	}

	if _, err = tmpFile.Write(data); err != nil {
		log.Fatalf("Failed to write to temp config file: %v", err)
	}
	if err = tmpFile.Close(); err != nil {
		log.Fatalf("Failed to close temp config file: %v", err)
	}

	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
}

func setupMockServer(t *testing.T) (*httptest.Server, string, *RequestCapture) {
	capture := &RequestCapture{}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Failed to read request body")
		}
		defer r.Body.Close()

		receivedRequest := HTTPRequestExpectation{
			URL:    r.URL.String(),
			Method: r.Method,
			Body:   string(body),
		}

		capture.Mutex.Lock()
		capture.ReceivedRequests = append(capture.ReceivedRequests, receivedRequest)
		capture.Mutex.Unlock()

		switch r.URL.Path {
		case "/400":
			w.WriteHeader(http.StatusBadRequest)
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))

	return mockServer, mockServer.URL, capture
}

func sleep(test TestCase) {
	retry := test.KonsumeConfig.Queues[0].Retry
	if retry == nil {
		time.Sleep(2 * time.Second)
		return
	}

	if retry.Enabled && retry.Interval > 0 && retry.Strategy != common.RetryStrategyExpo {
		time.Sleep(retry.Interval * time.Duration(retry.MaxRetries+1))
		return
	} else if retry.Enabled && retry.Interval > 0 && retry.Strategy == common.RetryStrategyExpo {
		time.Sleep(time.Duration(retry.MaxRetries*retry.MaxRetries) * retry.Interval)
		return
	}
}
