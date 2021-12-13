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
	"io/ioutil"
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

//func createDockerTarFile(dockerfile string) (io.Reader, error) {
//
//	var buf bytes.Buffer
//	tw := tar.NewWriter(&buf)
//	f, err := os.Open(dockerfile)
//	if err != nil {
//		println(dockerfile)
//		return nil, err
//	}
//	defer f.Close()
//
//	body, err := ioutil.ReadAll(f)
//	if err != nil {
//		return nil, err
//	}
//
//	hdr := &tar.Header{
//		Name: f.Name(),
//		Mode: 0600,
//		Size: int64(len(body)),
//	}
//
//	tw.WriteHeader(hdr)
//	_,err = tw.Write(body)
//	if err != nil {
//		return nil, err
//	}
//	return &buf, nil
//}

func (s ContainerService) BuildImage(opt model.BuildImageOption) error {
	ctx := context.Background()
	//dockerBuildContext, err := createDockerTarFile(opt.Path + opt.Dockerfile)
	//if err != nil {
	//	return err
	//}

	tar := new(archivex.TarFile)
	tar.Create(opt.Path)
	tar.AddAll(opt.Path, false)
	tar.Close()
	dockerBuildContext, err := os.Open(opt.Path + ".tar")
	defer dockerBuildContext.Close()
	resp, err := s.cli.ImageBuild(ctx, dockerBuildContext, types.ImageBuildOptions{
		Target:     opt.Tag,
		Version:    types.BuilderVersion(opt.Version),
		Dockerfile: opt.Path + opt.Dockerfile,
	})
	if err != nil {
		return err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
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
