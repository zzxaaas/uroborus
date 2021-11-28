package fx

import (
	"go.uber.org/fx"
	. "uroborus/service"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewProjectService),
	fx.Provide(NewGitService),
	fx.Provide(NewBaseImageService),
)
