package observability

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/DuongVu089x/interview/customer/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// WebSocketConnections tracks active WebSocket connections
	WebSocketConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "customer_service",
		Name:      "websocket_active_connections",
		Help:      "Number of active WebSocket connections",
	})

	// MessageProcessingDuration tracks message processing time
	MessageProcessingDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "customer_service",
		Name:      "message_processing_duration_seconds",
		Help:      "Time spent processing messages",
		Buckets:   prometheus.DefBuckets,
	}, []string{"topic", "status"})

	// NotificationsSent tracks number of notifications sent
	NotificationsSent = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "customer_service",
		Name:      "notifications_sent_total",
		Help:      "Total number of notifications sent",
	}, []string{"type"})

	// DatabaseOperationDuration tracks database operation duration
	DatabaseOperationDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "customer_service",
		Name:      "database_operation_duration_seconds",
		Help:      "Time spent on database operations",
		Buckets:   prometheus.DefBuckets,
	}, []string{"operation"})

	initOnce sync.Once
	registry *prometheus.Registry
)

// InitMetrics initializes Prometheus metrics collection
func InitMetrics(cfg *config.ObservabilityConfig) error {
	var initErr error
	initOnce.Do(func() {
		// Create a new registry
		registry = prometheus.NewRegistry()

		// Register standard collectors
		if err := registry.Register(collectors.NewGoCollector()); err != nil {
			initErr = fmt.Errorf("failed to register Go collector: %v", err)
			return
		}
		if err := registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
			initErr = fmt.Errorf("failed to register process collector: %v", err)
			return
		}

		// Register custom metrics
		collectors := []prometheus.Collector{
			WebSocketConnections,
			MessageProcessingDuration,
			NotificationsSent,
			DatabaseOperationDuration,
		}

		for _, collector := range collectors {
			if err := registry.Register(collector); err != nil {
				// If the collector is already registered, continue
				if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
					continue
				}
				initErr = fmt.Errorf("failed to register collector: %v", err)
				return
			}
		}

		// Start metrics server
		go func() {
			http.Handle(cfg.Prometheus.Endpoint, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
			if err := http.ListenAndServe(":"+cfg.Prometheus.Port, nil); err != nil {
				fmt.Printf("Error starting metrics server: %v\n", err)
			}
		}()
	})

	return initErr
}

// GetRegistry returns the Prometheus registry
func GetRegistry() *prometheus.Registry {
	return registry
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
