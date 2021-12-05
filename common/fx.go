package common

import "go.uber.org/fx"

// Module -
var Module = fx.Options(
	fx.Provide(NewDockerClient),
)
