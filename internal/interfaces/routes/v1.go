package routes

import (
	"github.com/Nezent/auth-service/internal/interfaces/handlers"
	"github.com/Nezent/auth-service/internal/interfaces/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type APIV1RoutesParams struct {
	fx.In

	Router      *chi.Mux
	Guard       *middleware.Guard
	UserHandler *handlers.UserHandler
}

type APIV1Routes struct {
	router      *chi.Mux
	guard       *middleware.Guard
	userHandler *handlers.UserHandler
}

func NewRoutes(params APIV1RoutesParams) *APIV1Routes {
	return &APIV1Routes{
		router:      params.Router,
		guard:       params.Guard,
		userHandler: params.UserHandler,
	}
}

func (r *APIV1Routes) Register() {
	r.router.Route("/v1", func(v1 chi.Router) {
		// guest routes
		v1.Route("/auth", func(noAuth chi.Router) {
			// noAuth.Post("/login", r.userHandler.Login)
			noAuth.Post("/register", r.userHandler.Register)
		})

		// auth routes (protected)
		v1.Route("/auth", func(auth chi.Router) {
			auth.Use(r.guard.Auth())
			// auth.Get("/me", r.userHandler.Me)
		})
	})
}
