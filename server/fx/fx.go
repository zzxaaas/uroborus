package fx

import (
	"go.uber.org/fx"
	. "uroborus/server"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewHealthServer),
	fx.Provide(NewProjectServer),
	fx.Provide(NewUserServer),
	fx.Provide(NewBaseImageServer),
	fx.Provide(NewDeployServer),
	fx.Provide(NewContainerServer),
	fx.Provide(NewGroupServer),
	fx.Provide(NewProjectCommentServer),
)
