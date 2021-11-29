package service

import (
	"github.com/docker/docker/client"
)

type DockerService struct {
	cli *client.Client
}

func NewDockerService(cli *client.Client) *DockerService {
	return &DockerService{
		cli: cli,
	}
}

//func (s DockerService) PullImage(image string) error {
//	ctx := context.Background()
//	reader,err := s.cli.ImagePull(ctx,image,types.ImagePullOptions{})
//	if err != nil {
//		return err
//	}
//	io.Copy(os.Stdout,reader)
//}
