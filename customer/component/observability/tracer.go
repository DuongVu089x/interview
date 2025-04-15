package observability

import (
	"fmt"

	"github.com/DuongVu089x/interview/customer/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes the Jaeger tracer
func InitTracer(cfg *config.ObservabilityConfig) (*tracesdk.TracerProvider, error) {
	if !cfg.Jaeger.Enabled {
		return nil, nil
	}

	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(cfg.Jaeger.AgentHost),
		jaeger.WithAgentPort(cfg.Jaeger.AgentPort),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Jaeger.ServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// GetTracer returns a named tracer
func GetTracer(name string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name)
}
