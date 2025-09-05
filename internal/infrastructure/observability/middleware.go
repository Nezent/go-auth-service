package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// PrometheusMiddleware wraps the metrics for HTTP middleware
type PrometheusMiddleware struct {
	metrics *Metrics
}

// NewPrometheusMiddleware creates a new Prometheus middleware
func NewPrometheusMiddleware(metrics *Metrics) *PrometheusMiddleware {
	return &PrometheusMiddleware{
		metrics: metrics,
	}
}

// Middleware returns the HTTP middleware function for Prometheus metrics
func (pm *PrometheusMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Increment in-flight requests
			pm.metrics.IncHTTPRequestsInFlight()
			defer pm.metrics.DecHTTPRequestsInFlight()

			// Create a response wrapper to capture status code and response size
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process the request
			next.ServeHTTP(ww, r)

			// Calculate duration
			duration := time.Since(start)

			// Get route pattern for better endpoint grouping
			endpoint := r.URL.Path
			if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
				if routeContext.RoutePattern() != "" {
					endpoint = routeContext.RoutePattern()
				}
			}

			// Record metrics
			pm.metrics.RecordHTTPRequest(
				r.Method,
				endpoint,
				strconv.Itoa(ww.Status()),
				duration,
				float64(ww.BytesWritten()),
			)
		})
	}
}

// RecoveryMiddleware returns a middleware that records panic recoveries
func (pm *PrometheusMiddleware) RecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					pm.metrics.RecordPanicRecovery()
					pm.metrics.RecordError("http", "panic")

					// Re-panic to let the default recovery middleware handle it
					panic(rec)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
