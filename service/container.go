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
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/jhoonb/archivex"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"uroborus/common/kafka"
	"uroborus/model"
)

type ContainerService struct {
	cli            *client.Client
	projectService *ProjectService
	kafkaCli       *kafka.Client
	upGrader       websocket.Upgrader
	redisCli       *redis.Client
}

func NewContainerService(
	cli *client.Client,
	kafkaCli *kafka.Client,
	projectService *ProjectService,
	redisCli *redis.Client,
) *ContainerService {
	return &ContainerService{
		cli:            cli,
		kafkaCli:       kafkaCli,
		projectService: projectService,
		upGrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		redisCli: redisCli,
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

func (s ContainerService) getCpuUsage(key string) (uint64, uint64) {
	preCpuUsage, _ := s.redisCli.Get(key + model.KeyCtnPreCpuSuffix).Uint64()
	preSysUsage, _ := s.redisCli.Get(key + model.KeyCtnPreSysSuffix).Uint64()
	return preCpuUsage, preSysUsage
}

func (s ContainerService) setCpuUsage(key string, cpuUsage, sysUsage uint64) {
	s.redisCli.Set(key+model.KeyCtnPreCpuSuffix, cpuUsage, -1)
	s.redisCli.Set(key+model.KeyCtnPreSysSuffix, sysUsage, -1)
	return
}

func (s ContainerService) GetAll(user string) (error, []model.GetContainerStatsResp) {
	ctx := context.Background()
	resp := make([]model.GetContainerStatsResp, 0)
	projects, err := s.projectService.Find(model.Project{UserName: user})
	if err != nil {
		return err, nil
	}
	for _, p := range projects {
		stats := types.StatsJSON{}
		if p.Container == "" {
			resp = append(resp, model.GetContainerStatsResp{Name: p.Name})
			continue
		}
		reader, err := s.cli.ContainerStatsOneShot(ctx, p.Container)
		if err != nil {
			return err, nil
		}

		body, err := ioutil.ReadAll(reader.Body)
		json.Unmarshal(body, &stats)
		tmp := model.GetContainerStatsResp{
			Name:     stats.Name[1:],
			ID:       stats.ID[:12],
			MemLimit: stats.MemoryStats.Limit,
			MemUsage: stats.MemoryStats.Usage,
			NETIn:    stats.Networks["eth0"].RxBytes,
			NETOut:   stats.Networks["eth0"].TxBytes,
		}
		key := fmt.Sprintf("%s:%s", model.RedisKeyCtnPrefix, tmp.ID)
		preCpuUsage, preSysUsage := s.getCpuUsage(key)
		if preCpuUsage != 0 || preSysUsage != 0 {
			tmp.CPUUsage = calculateCPUPercentUnix(0, 0, &stats)
		}
		s.setCpuUsage(key, stats.CPUStats.CPUUsage.TotalUsage, stats.CPUStats.SystemUsage)
		tmp.MemPercent = float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit) * 100
		resp = append(resp, tmp)

	}
	return nil, resp
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 1000.0
	}
	return cpuPercent
}

func (s ContainerService) Terminal(conn *websocket.Conn, req model.ConnectContainerReq) error {
	// 执行exec，获取到容器终端的连接
	hr, err := s.exec(req.ID)
	if err != nil {
		return err
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	// 转发输入/输出至websocket
	go func() {
		wsWriterCopy(hr.Conn, conn)
	}()
	wsReaderCopy(conn, hr.Conn)
	return nil
}

func (s ContainerService) exec(container string) (hr types.HijackedResponse, err error) {
	ctx := context.Background()
	// 执行/bin/bash命令
	ir, err := s.cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	})
	if err != nil {
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = s.cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return
	}
	return
}

// 将终端的输出转发到前端
func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

// 将前端的输入转发到终端
func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}
