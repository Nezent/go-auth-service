package router

import (
	"net/http"
	"time"

	"github.com/Nezent/auth-service/internal/infrastructure/config"
	"github.com/Nezent/auth-service/internal/infrastructure/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(prometheusMiddleware *observability.PrometheusMiddleware) *chi.Mux {
	cfg := config.NewConfig()
	router := chi.NewRouter()

	// Basic middleware (order matters!)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)

	// Custom recovery middleware that records panics
	router.Use(prometheusMiddleware.RecoveryMiddleware())
	router.Use(middleware.Recoverer)

	// Prometheus metrics middleware (early in chain)
	router.Use(prometheusMiddleware.Middleware())

	// Request size limiting (prevent large payloads)
	router.Use(middleware.RequestSize(1024 * 1024)) // 1MB limit

	// Compression middleware (before security headers)
	router.Use(middleware.Compress(5)) // gzip compression

	// Security headers middleware
	router.Use(securityHeadersMiddleware(cfg))

	// CORS middleware
	if cfg.CORS.Enabled {
		router.Use(corsMiddleware(cfg))
	}

	// Rate limiting middleware
	if cfg.RateLimit.Enabled {
		router.Use(httprate.LimitByIP(cfg.RateLimit.RequestsPerMinute, time.Minute))
	}

	// Timeout middleware
	router.Use(middleware.Timeout(30 * time.Second))

	// Content type middleware for API responses
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Remove trailing slashes
	router.Use(middleware.RedirectSlashes)

	// Profiler middleware (only in development)
	if cfg.Service.Debug && cfg.Service.Env == "development" {
		router.Mount("/debug", middleware.Profiler())
	}

	// Heartbeat endpoint for health checks
	router.Use(middleware.Heartbeat("/health"))

	// Readiness probe endpoint for Kubernetes
	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready","service":"` + cfg.Service.Name + `","version":"` + cfg.Service.Version + `"}`))
	})

	// Liveness probe endpoint for Kubernetes
	router.Get("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"alive"}`))
	})

	// Metrics endpoint
	if cfg.Monitoring.Enabled {
		router.Handle(cfg.Monitoring.PrometheusEndpoint, promhttp.Handler())
	}

	return router
}

// corsMiddleware configures CORS based on the application configuration
func corsMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.CORS.ExposedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	})
}

// securityHeadersMiddleware adds various security headers
func securityHeadersMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			if cfg.Security.ContentTypeNosniff {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}

			// Clickjacking protection
			if cfg.Security.FrameOptions != "" {
				w.Header().Set("X-Frame-Options", cfg.Security.FrameOptions)
			}

			// XSS protection
			if cfg.Security.XSSProtection != "" {
				w.Header().Set("X-XSS-Protection", cfg.Security.XSSProtection)
			}

			// Content Security Policy
			if cfg.Security.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", cfg.Security.ContentSecurityPolicy)
			}

			// Referrer Policy
			if cfg.Security.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", cfg.Security.ReferrerPolicy)
			}

			// HSTS (HTTP Strict Transport Security) for HTTPS
			// if r.TLS != nil && cfg.Security.HSTSMaxAge > 0 {
			// 	hstsValue := fmt.Sprintf("max-age=%d", cfg.Security.HSTSMaxAge)
			// 	if cfg.Security.HSTSIncludeSubdomains {
			// 		hstsValue += "; includeSubDomains"
			// 	}
			// 	w.Header().Set("Strict-Transport-Security", hstsValue)
			// }

			// Remove server information
			w.Header().Set("Server", "")

			next.ServeHTTP(w, r)
		})
	}
}
