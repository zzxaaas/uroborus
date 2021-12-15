package common

import (
	"github.com/docker/docker/client"
)

func NewDockerClient() *client.Client {
	//cli, err := client.NewClientWithOpts(
	//	client.WithHost(viper.GetString("docker.host")),
	//)
	cli, err := client.NewClientWithOpts(
		client.WithHTTPHeaders(map[string]string{"Content-Type": "application/tar"}),
	)
	if err != nil {
		panic(err.Error())
	}
	return cli
}
