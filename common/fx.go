package common

import (
	"go.uber.org/fx"
	. "uroborus/common/docker"
	. "uroborus/common/kafka"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewDockerClient),
	fx.Provide(NewKafkaClient),
)
