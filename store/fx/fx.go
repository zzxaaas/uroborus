package fx

import (
	"go.uber.org/fx"
	. "uroborus/store"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewPgDB),
	fx.Provide(NewUserStore),
	fx.Provide(NewProjectStore),
	fx.Provide(NewBaseImageStore),
	fx.Provide(NewDeployHistoryStore),
)
