package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// contextKey type for context keys
type contextKey string

const (
	correlationIDKey contextKey = "correlation_id"
	requestStartKey  contextKey = "request_start"
)

// Prometheus metrics
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served",
		},
	)
)

// CorrelationIDMiddleware adds a correlation ID to each request
// If the X-Correlation-ID header is present, its value is used
// Otherwise, a new UUID is generated
// Chi-compatible middleware
func CorrelationIDMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Try to get correlation ID from header
			correlationID := r.Header.Get("X-Correlation-ID")
			if correlationID == "" {
				// Generate a new UUID
				correlationID = uuid.New().String()
			}

			// Add correlation ID to context
			ctx := context.WithValue(r.Context(), correlationIDKey, correlationID)

			// Add correlation ID to response header
			w.Header().Set("X-Correlation-ID", correlationID)

			// Log request start
			logger.Info(ctx, "request started", map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			})

			// Pass control to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// GetCorrelationID extracts the correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value(correlationIDKey).(string); ok {
		return correlationID
	}
	return ""
}

// MetricsMiddleware collects request metrics using Prometheus
// Chi-compatible middleware
func MetricsMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Increment in-flight requests counter
			httpRequestsInFlight.Inc()
			defer httpRequestsInFlight.Dec()

			// Add start time to context
			ctx := context.WithValue(r.Context(), requestStartKey, start)

			// Wrap ResponseWriter to capture status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Pass control to next handler
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Calculate duration
			duration := time.Since(start)

			// Update Prometheus metrics
			httpRequestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(rw.statusCode),
			).Inc()

			httpRequestDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
			).Observe(duration.Seconds())

			// Log request completion
			logger.Info(ctx, "request completed", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status":      rw.statusCode,
				"duration_ms": duration.Milliseconds(),
			})
		}
		return http.HandlerFunc(fn)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

// Write captures data writing
func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(data)
}

// RecoveryMiddleware recovers from panics
// Chi-compatible middleware
func RecoveryMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(r.Context(), "panic recovered", map[string]interface{}{
						"error":  err,
						"method": r.Method,
						"path":   r.URL.Path,
					})

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal Server Error","message":"an unexpected error occurred"}`))
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
