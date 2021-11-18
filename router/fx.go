package router

import (
	"go.uber.org/fx"
	"uroborus/router/middleware"
)

// Module -
var Module = fx.Options(
	middleware.Module,
	fx.Provide(NewRouter),
	fx.Provide(NewHTTPServer),
)
