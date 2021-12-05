package common

import (
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

func NewDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(
		client.WithHost(viper.GetString("docker.host")),
	)
	if err != nil {
		panic(err.Error())
	}
	return cli
}
