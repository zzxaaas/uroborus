package fx

import (
	"go.uber.org/fx"
	. "uroborus/store"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewPgDB),
	fx.Provide(NewUserStore),
)
