package middleware

import (
	"context"
	"net/http"

	"github.com/Nezent/auth-service/internal/constants"
	"github.com/Nezent/auth-service/internal/infrastructure/persistence"
	"github.com/Nezent/auth-service/pkg/response"
)

type GuardOptions struct {
	Permission          constants.Permission
	RequireSubscription bool
	Handler             http.HandlerFunc
}

type Guard struct {
	db *persistence.Database
}

func NewGuard(db *persistence.Database) *Guard {
	return &Guard{db: db}
}

// Permission check middleware
func (g *Guard) Handler(options *GuardOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: get user permissions from context/session/db
		userPerm := getUserPermission(r.Context()) // implement this function

		if userPerm != options.Permission {
			response.WriteError(w, "Permission denied", http.StatusForbidden)
			return
		}

		if options.RequireSubscription && !hasActiveSubscription(r.Context()) { // implement this function
			response.WriteError(w, "Subscription required", http.StatusPaymentRequired)
			return
		}

		options.Handler(w, r)
	}
}

// Simple auth middleware (for demonstration)
func (g *Guard) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Example: check authentication here
			if !isAuthenticated(r.Context()) {
				response.WriteError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// --- Helper functions (implement as needed) ---

func getUserPermission(ctx context.Context) constants.Permission {
	// TODO: Extract user permission from context/session/database
	return constants.OrgViewHealth // example
}

func hasActiveSubscription(ctx context.Context) bool {
	// TODO: Check subscription status from context/session/database
	return true // example
}

func isAuthenticated(ctx context.Context) bool {
	// TODO: Check authentication from context/session/database
	return true // example
}
