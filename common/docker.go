package common

import (
	"github.com/docker/docker/client"
)

func NewDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err.Error())
	}
	return cli
}
