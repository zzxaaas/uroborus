package fx

import (
	"go.uber.org/fx"
	. "uroborus/server"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewHealthServer),
	fx.Provide(NewProjectServer),
)
