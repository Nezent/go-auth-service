package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics for the auth service
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPResponseSize     *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// Authentication metrics
	AuthAttempts        *prometheus.CounterVec
	AuthTokensGenerated *prometheus.CounterVec
	AuthTokensValidated *prometheus.CounterVec
	AuthFailures        *prometheus.CounterVec

	// Business metrics
	UsersRegistered prometheus.Counter
	UsersActive     prometheus.Gauge
	LoginAttempts   *prometheus.CounterVec
	PasswordResets  prometheus.Counter

	// Database metrics
	DatabaseConnections   prometheus.Gauge
	DatabaseQueries       *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec

	// Cache metrics
	CacheHits       *prometheus.CounterVec
	CacheMisses     *prometheus.CounterVec
	CacheOperations *prometheus.CounterVec

	// System metrics
	ErrorsTotal     *prometheus.CounterVec
	PanicRecoveries prometheus.Counter
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "auth_service",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "auth_service",
				Subsystem: "http",
				Name:      "response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"method", "endpoint"},
		),

		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "auth_service",
				Subsystem: "http",
				Name:      "requests_in_flight",
				Help:      "Current number of HTTP requests being processed",
			},
		),

		// Authentication metrics
		AuthAttempts: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "auth",
				Name:      "attempts_total",
				Help:      "Total number of authentication attempts",
			},
			[]string{"type", "status"},
		),

		AuthTokensGenerated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "auth",
				Name:      "tokens_generated_total",
				Help:      "Total number of tokens generated",
			},
			[]string{"type"},
		),

		AuthTokensValidated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "auth",
				Name:      "tokens_validated_total",
				Help:      "Total number of tokens validated",
			},
			[]string{"type", "status"},
		),

		AuthFailures: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "auth",
				Name:      "failures_total",
				Help:      "Total number of authentication failures",
			},
			[]string{"reason"},
		),

		// Business metrics
		UsersRegistered: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "users",
				Name:      "registered_total",
				Help:      "Total number of users registered",
			},
		),

		UsersActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "auth_service",
				Subsystem: "users",
				Name:      "active_current",
				Help:      "Current number of active users",
			},
		),

		LoginAttempts: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "users",
				Name:      "login_attempts_total",
				Help:      "Total number of login attempts",
			},
			[]string{"status"},
		),

		PasswordResets: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "users",
				Name:      "password_resets_total",
				Help:      "Total number of password reset requests",
			},
		),

		// Database metrics
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "auth_service",
				Subsystem: "database",
				Name:      "connections_current",
				Help:      "Current number of database connections",
			},
		),

		DatabaseQueries: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "database",
				Name:      "queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),

		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "auth_service",
				Subsystem: "database",
				Name:      "query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
			},
			[]string{"operation", "table"},
		),

		// Cache metrics
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "cache",
				Name:      "hits_total",
				Help:      "Total number of cache hits",
			},
			[]string{"cache_name"},
		),

		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "cache",
				Name:      "misses_total",
				Help:      "Total number of cache misses",
			},
			[]string{"cache_name"},
		),

		CacheOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "cache",
				Name:      "operations_total",
				Help:      "Total number of cache operations",
			},
			[]string{"operation", "cache_name", "status"},
		),

		// System metrics
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "system",
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"component", "error_type"},
		),

		PanicRecoveries: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: "auth_service",
				Subsystem: "system",
				Name:      "panic_recoveries_total",
				Help:      "Total number of panic recoveries",
			},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration, responseSize float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.HTTPResponseSize.WithLabelValues(method, endpoint).Observe(responseSize)
}

// IncHTTPRequestsInFlight increments in-flight requests
func (m *Metrics) IncHTTPRequestsInFlight() {
	m.HTTPRequestsInFlight.Inc()
}

// DecHTTPRequestsInFlight decrements in-flight requests
func (m *Metrics) DecHTTPRequestsInFlight() {
	m.HTTPRequestsInFlight.Dec()
}

// RecordAuthAttempt records authentication attempt
func (m *Metrics) RecordAuthAttempt(authType, status string) {
	m.AuthAttempts.WithLabelValues(authType, status).Inc()
}

// RecordTokenGenerated records token generation
func (m *Metrics) RecordTokenGenerated(tokenType string) {
	m.AuthTokensGenerated.WithLabelValues(tokenType).Inc()
}

// RecordTokenValidated records token validation
func (m *Metrics) RecordTokenValidated(tokenType, status string) {
	m.AuthTokensValidated.WithLabelValues(tokenType, status).Inc()
}

// RecordAuthFailure records authentication failure
func (m *Metrics) RecordAuthFailure(reason string) {
	m.AuthFailures.WithLabelValues(reason).Inc()
}

// RecordUserRegistration records user registration
func (m *Metrics) RecordUserRegistration() {
	m.UsersRegistered.Inc()
}

// SetActiveUsers sets the current number of active users
func (m *Metrics) SetActiveUsers(count float64) {
	m.UsersActive.Set(count)
}

// RecordLoginAttempt records login attempt
func (m *Metrics) RecordLoginAttempt(status string) {
	m.LoginAttempts.WithLabelValues(status).Inc()
}

// RecordPasswordReset records password reset
func (m *Metrics) RecordPasswordReset() {
	m.PasswordResets.Inc()
}

// RecordDatabaseQuery records database query metrics
func (m *Metrics) RecordDatabaseQuery(operation, table, status string, duration time.Duration) {
	m.DatabaseQueries.WithLabelValues(operation, table, status).Inc()
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// SetDatabaseConnections sets the current number of database connections
func (m *Metrics) SetDatabaseConnections(count float64) {
	m.DatabaseConnections.Set(count)
}

// RecordCacheHit records cache hit
func (m *Metrics) RecordCacheHit(cacheName string) {
	m.CacheHits.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss records cache miss
func (m *Metrics) RecordCacheMiss(cacheName string) {
	m.CacheMisses.WithLabelValues(cacheName).Inc()
}

// RecordCacheOperation records cache operation
func (m *Metrics) RecordCacheOperation(operation, cacheName, status string) {
	m.CacheOperations.WithLabelValues(operation, cacheName, status).Inc()
}

// RecordError records system error
func (m *Metrics) RecordError(component, errorType string) {
	m.ErrorsTotal.WithLabelValues(component, errorType).Inc()
}

// RecordPanicRecovery records panic recovery
func (m *Metrics) RecordPanicRecovery() {
	m.PanicRecoveries.Inc()
}
