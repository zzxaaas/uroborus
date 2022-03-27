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
	fx.Provide(NewDockerService),
	fx.Provide(NewContainerService),
	fx.Provide(NewDeployService),
	fx.Provide(NewDeployHistoryService),
	fx.Provide(NewDeployLogService),
	fx.Provide(NewProxyService),
	fx.Provide(NewGroupService),
)
