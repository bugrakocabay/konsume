package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bugrakocabay/konsume/pkg/config"
)

var (
	MessagesConsumed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konsume_messages_consumed_total",
		Help: "Total number of messages consumed",
	})

	HttpRequestsMade = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konsume_http_requests_made_total",
		Help: "Total number of HTTP requests made",
	})

	HttpRequestsSucceeded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konsume_http_requests_succeeded_total",
		Help: "Total number of HTTP requests succeeded",
	})

	HttpRequestsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "konsume_http_requests_failed_total",
		Help: "Total number of HTTP requests failed",
	})
)

// InitMetrics initializes the metrics endpoint for prometheus with custom metrics
func InitMetrics(cfg *config.MetricsConfig) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(MessagesConsumed)
	registry.MustRegister(HttpRequestsMade)
	registry.MustRegister(HttpRequestsSucceeded)
	registry.MustRegister(HttpRequestsFailed)
	registry.MustRegister(collectors.NewBuildInfoCollector())
	registry.MustRegister(collectors.NewGoCollector())

	go func() {
		http.Handle(cfg.Path, promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
	}()
}
