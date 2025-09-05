package routes

import "go.uber.org/fx"

// Provide routes in Uber FX module
var Module = fx.Module("routes",
	fx.Provide(NewRoutes),
	fx.Invoke(func(r *APIV1Routes) {
		r.Register()
	}),
)
