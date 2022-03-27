package common

import (
	"go.uber.org/fx"
	. "uroborus/common/docker"
	. "uroborus/common/kafka"
	. "uroborus/common/redis"
)

// Module -
var Module = fx.Options(
	fx.Provide(NewDockerClient),
	fx.Provide(NewKafkaClient),
	fx.Provide(NewRedisClient),
)
