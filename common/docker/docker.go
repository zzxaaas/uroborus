package docker

import (
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

func NewDockerClient() *client.Client {
	//cli, err := client.NewClientWithOpts(
	//	client.WithHost(viper.GetString("docker.host")),
	//)
	cli, err := client.NewClientWithOpts(
		client.WithHost(viper.GetString("docker.host")),
		client.WithHTTPHeaders(map[string]string{"Content-Type": "application/tar"}),
	)
	if err != nil {
		panic(err.Error())
	}
	return cli
}
