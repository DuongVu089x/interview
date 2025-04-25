package middleware

import (
	"time"

	"github.com/DuongVu089x/interview/order/component/appctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ObservabilityMiddleware adds logging, tracing, and metrics to each request
func ObservabilityMiddleware(appCtx appctx.AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := req.URL.Path
			method := req.Method
			requestId := uuid.New().String()

			// Get logger with request context
			logger := appCtx.GetLogger().With(
				zap.String("path", path),
				zap.String("method", method),
				zap.String("request_id", requestId),
			)

			// Start tracing
			ctx := req.Context()
			tracer := appCtx.GetTracer()
			ctx, span := tracer.Start(ctx, "http.request")
			defer span.End()

			// Add common trace attributes
			span.SetAttributes(
				attribute.String("http.path", path),
				attribute.String("http.method", method),
				attribute.String("request_id", requestId),
			)

			// Add logger and tracer to context
			c.Set("logger", logger)
			c.Set("span", span)

			// Update request context
			c.SetRequest(req.WithContext(ctx))

			// Start timing
			start := time.Now()

			// Process request
			err := next(c)

			// Calculate request duration
			duration := time.Since(start).Seconds()
			status := c.Response().Status

			// Add response attributes to span
			span.SetAttributes(
				attribute.Int("http.status_code", status),
				attribute.Float64("duration_seconds", duration),
			)

			// Log request completion
			logger = logger.With(
				zap.Int("status", status),
				zap.Float64("duration_ms", duration*1000),
			)

			if err != nil {
				logger.Error("Request failed", zap.Error(err))
				span.RecordError(err)
			} else {
				logger.Info("Request completed successfully")
			}

			return err
		}
	}
}

// GetRequestLogger gets the logger from context
func GetRequestLogger(c echo.Context) *zap.Logger {
	if logger, ok := c.Get("logger").(*zap.Logger); ok {
		return logger
	}
	// Fallback to default logger if not found in context
	return zap.L()
}

// GetRequestSpan gets the trace span from context
func GetRequestSpan(c echo.Context) trace.Span {
	if span, ok := c.Get("span").(trace.Span); ok {
		return span
	}
	return nil
}
