package service

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService(cli *client.Client) *DockerService {
	return &DockerService{
		cli: cli,
	}
}

func (s DockerService) PullImage(image string) error {
	ctx := context.Background()
	reader, err := s.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)
	return nil
}

func (s DockerService) StartContainer(image, containerName string) {
	ctx := context.Background()

	resp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, nil, nil, containerName)
	if err != nil {
		panic(err)
	}

	if err := s.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

}
