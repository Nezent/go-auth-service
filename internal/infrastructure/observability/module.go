package observability

import (
	"go.uber.org/fx"
)

// Module provides observability dependencies for the application
var Module = fx.Options(
	fx.Provide(
		NewMetrics,
		NewPrometheusMiddleware,
	),
)
