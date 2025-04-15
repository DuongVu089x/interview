package observability

import (
	"fmt"
	"net/http"

	"github.com/DuongVu089x/interview/customer/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// WebSocketConnections tracks active WebSocket connections
	WebSocketConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "websocket_active_connections",
		Help: "Number of active WebSocket connections",
	})

	// MessageProcessingDuration tracks message processing time
	MessageProcessingDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "message_processing_duration_seconds",
		Help:    "Time spent processing messages",
		Buckets: prometheus.DefBuckets,
	}, []string{"topic", "status"})

	// NotificationsSent tracks number of notifications sent
	NotificationsSent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "notifications_sent_total",
		Help: "Total number of notifications sent",
	}, []string{"type"})

	// DatabaseOperationDuration tracks database operation duration
	DatabaseOperationDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_operation_duration_seconds",
		Help:    "Time spent on database operations",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})
)

// InitMetrics initializes Prometheus metrics collection
func InitMetrics(cfg *config.ObservabilityConfig) error {
	// Register standard collectors
	prometheus.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Register custom metrics
	prometheus.MustRegister(
		WebSocketConnections,
		MessageProcessingDuration,
		NotificationsSent,
		DatabaseOperationDuration,
	)

	// Start metrics server
	go func() {
		http.Handle(cfg.Prometheus.Endpoint, promhttp.Handler())
		if err := http.ListenAndServe(":"+cfg.Prometheus.Port, nil); err != nil {
			fmt.Printf("Error starting metrics server: %v\n", err)
		}
	}()

	return nil
}

// RecordWebSocketConnection records WebSocket connection metrics
func RecordWebSocketConnection(delta float64) {
	WebSocketConnections.Add(delta)
}

// RecordMessageProcessing records message processing duration
func RecordMessageProcessing(topic string, status string, duration float64) {
	MessageProcessingDuration.WithLabelValues(topic, status).Observe(duration)
}

// RecordNotificationSent records a sent notification
func RecordNotificationSent(notificationType string) {
	NotificationsSent.WithLabelValues(notificationType).Inc()
}

// RecordDatabaseOperation records database operation duration
func RecordDatabaseOperation(operation string, duration float64) {
	DatabaseOperationDuration.WithLabelValues(operation).Observe(duration)
}
