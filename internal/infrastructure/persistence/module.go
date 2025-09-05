package persistence

import "go.uber.org/fx"

var Module = fx.Module(
	"database",
	fx.Provide(
		NewDatabase,
	),
)
