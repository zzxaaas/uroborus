package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jhoonb/archivex"
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

func (s ContainerService) BuildImage(opt model.BuildImageOption) error {
	ctx := context.Background()

	tar := new(archivex.TarFile)
	tar.Create(opt.Path[:len(opt.Path)-1] + ".tar")
	tar.AddAll(opt.Path[:len(opt.Path)-1], false)
	tar.Close()
	dockerBuildContext, err := os.Open(opt.Path[:len(opt.Path)-1] + ".tar")
	defer dockerBuildContext.Close()
	resp, err := s.cli.ImageBuild(ctx, dockerBuildContext, types.ImageBuildOptions{
		Tags:       []string{opt.Tag},
		Dockerfile: opt.Dockerfile,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	//response, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(string(response))
	return nil
}

func (s ContainerService) RemoveContainer(contianerID string) error {
	ctx := context.Background()
	s.cli.ContainerStop(ctx, contianerID, nil)
	return s.cli.ContainerRemove(ctx, contianerID, types.ContainerRemoveOptions{})
}

func (s ContainerService) StartContainerWithOption(opt model.ContainerOption) (string, error) {
	ctx := context.Background()
	fmt.Println(opt)
	if opt.NeedPull {
		out, err := s.cli.ImagePull(ctx, opt.Image, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	}

	_, portMap, err := nat.ParsePortSpecs([]string{fmt.Sprintf("%s:%s", opt.Port, opt.ProtoPort)})
	if err != nil {
		return "", err
	}

	ctrConfig := container.Config{
		Image: opt.Image,
	}
	if len(opt.Env) > 0 && opt.Env[0] != "" {
		ctrConfig.Env = opt.Env
	}

	resp, err := s.cli.ContainerCreate(ctx, &ctrConfig, &container.HostConfig{
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
