package service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/websocket"
	"github.com/jhoonb/archivex"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"uroborus/common/kafka"
	"uroborus/model"
)

type ContainerService struct {
	cli      *client.Client
	kafkaCli *kafka.Client
	upGrader websocket.Upgrader
}

func NewContainerService(cli *client.Client, kafkaCli *kafka.Client) *ContainerService {
	return &ContainerService{
		cli:      cli,
		kafkaCli: kafkaCli,
		upGrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
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
		Tags: []string{opt.Tag},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	s.ReaderToKafka(resp.Body, int(opt.DeployID), model.DEPLOY_STEP_BUILD)
	return nil
}

func (s ContainerService) GetContainerLog(contianerID string) (io.Reader, error) {
	logs, err := s.cli.ContainerLogs(context.Background(), contianerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s ContainerService) RemoveContainer(contianerID string) error {
	ctx := context.Background()
	s.cli.ContainerStop(ctx, contianerID, nil)
	s.cli.ContainerRemove(ctx, contianerID, types.ContainerRemoveOptions{})
	return nil
}

func (s ContainerService) StartContainerWithOption(opt model.ContainerOption) (string, error) {
	ctx := context.Background()
	if opt.NeedPull {
		out, err := s.cli.ImagePull(ctx, opt.Image, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer out.Close()
		s.ReaderToKafka(out, int(opt.DeployID), model.DEPLOY_STEP_BUILD)
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
	topic := strconv.Itoa(int(opt.DeployID))
	resp, err := s.cli.ContainerCreate(ctx, &ctrConfig, &container.HostConfig{
		PortBindings: portMap,
	}, nil, nil, opt.Name)
	if err != nil {
		s.sendMessage(topic, err.Error(), model.DEPLOY_STEP_RUN)
		return "", err
	}
	s.sendMessage(topic, fmt.Sprintf("create container {%s} success\n", resp.ID), model.DEPLOY_STEP_RUN)
	if err := s.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	s.sendMessage(topic, "start container success\n", model.DEPLOY_STEP_RUN)
	return resp.ID, nil
}

func (s ContainerService) sendMessage(topic, message string, key int32) {
	s.kafkaCli.SendLog(kafka.PackMsg(topic, message, key))
}

func (s ContainerService) ReaderToKafka(reader io.Reader, deployID, step int) {
	r := bufio.NewReader(reader)
	for {
		log, err := r.ReadString('\n')
		if err != nil {
			logrus.Error(err)
			break
		}
		stream := model.LogStream{}
		json.Unmarshal([]byte(log), &stream)
		if strings.Contains(log, "status") {
			log = stream.Status + "\n"
		} else {
			log = stream.Stream
		}
		s.sendMessage(strconv.Itoa(deployID), log, int32(step))

	}
}
