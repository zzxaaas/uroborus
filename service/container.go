package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"uroborus/model"
)

type ContainerService struct {
	cli *client.Client
}

func NewContainerService(cli *client.Client) *ContainerService {
	return &ContainerService{
		cli: cli,
	}
}

func (s ContainerService) StartContainerWithDockerfile() {

}

func (s ContainerService) RemoveContainer(contianerID string) error {
	ctx := context.Background()
	s.cli.ContainerStop(ctx, contianerID, nil)
	return s.cli.ContainerRemove(ctx, contianerID, types.ContainerRemoveOptions{})
}

func (s ContainerService) StartContainerWithOption(opt model.ContainerOption) (string, error) {
	ctx := context.Background()

	out, err := s.cli.ImagePull(ctx, opt.Image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	_, portMap, err := nat.ParsePortSpecs([]string{fmt.Sprintf("%s:%s", opt.Port, opt.ProtoPort)})
	if err != nil {
		return "", err
	}
	resp, err := s.cli.ContainerCreate(ctx, &container.Config{
		Image: opt.Image,
		Env:   opt.Env,
	}, &container.HostConfig{
		PortBindings: portMap,
	}, nil, nil, opt.Name)
	if err != nil {
		return "", err
	}

	if err := s.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return resp.ID, nil
}
