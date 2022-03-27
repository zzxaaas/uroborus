package service

import (
	"bufio"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"uroborus/common"
	"uroborus/common/kafka"
	"uroborus/model"
)

type DeployLogService struct {
	deployHistoryService *DeployHistoryService
	projectService       *ProjectService
	gitService           *GitService
	containerService     *ContainerService
	kafkaCli             *kafka.Client
}

func NewDeployLogService(deployHistoryService *DeployHistoryService, projectService *ProjectService, gitService *GitService, containerService *ContainerService, kafkaCli *kafka.Client) *DeployLogService {
	return &DeployLogService{
		deployHistoryService: deployHistoryService,
		projectService:       projectService,
		gitService:           gitService,
		containerService:     containerService,
		kafkaCli:             kafkaCli,
	}
}

func (s DeployLogService) GetLog(conn *websocket.Conn, body *model.DeployHistory) error {
	if err := s.deployHistoryService.Get(body); err != nil {
		return err
	}
	project := &model.Project{Model: model.Model{ID: body.OriginId}}
	if _, err := s.projectService.Get(project); err != nil {
		return err
	}
	logPath := s.getLogPath(project, body.ID)
	//step, err := s.GetLogFromFile(conn, logPath)
	//if err != nil {
	//	return err
	//}
	step := 0
	if err := s.GetLogFromKafka(int(body.ID), step, logPath, conn); err != nil {
		return err
	}
	return nil
}

func (s DeployLogService) GetRunningLog(conn *websocket.Conn, body *model.DeployHistory) error {
	if err := s.deployHistoryService.Get(body); err != nil {
		return err
	}
	project := &model.Project{Model: model.Model{ID: body.OriginId}}
	if _, err := s.projectService.Get(project); err != nil {
		return err
	}
	s.GetContainerLog(conn, project.Container)
	return nil
}

func (s DeployLogService) GetLogFromKafka(deployID, step int, logPath string, conn *websocket.Conn) error {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	consumer, err := sarama.NewConsumer(viper.GetStringSlice("kafka.addrs"), config)
	if err != nil {
		return err
	}
	defer consumer.Close()
	topic := strconv.Itoa(deployID)
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, p := range partitions {
		pc, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
		defer pc.AsyncClose()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(pc sarama.PartitionConsumer, part int32) {
			defer wg.Done()
			for m := range pc.Messages() {
				if strconv.Itoa(step) > string(m.Key) {
					continue
				}
				value := string(m.Value)
				// todo:特殊处理未知字符导致乱码
				if strings.Contains(value, "go:") {
					value = string(m.Value[5 : len(m.Value)-5])
					value += "\n"
				}
				info := fmt.Sprintf("%s<->%s", string(m.Key), value)
				err = conn.WriteJSON(info)
				//s.SaveToFile(logPath, string(m.Key), m.Value)
				if err != nil {
					logrus.Error("err")
				}
			}
		}(pc, p)
		wg.Wait()
	}
	return nil
}

func (s DeployLogService) GetLogFromFile(conn *websocket.Conn, logPath string) (int, error) {
	lastStep := 0
	if has, _ := common.PathExists(logPath); has {
		logFiles, err := ioutil.ReadDir(logPath)
		if err != nil {
			return lastStep, err
		}
		for _, logFile := range logFiles {
			lastStep, err = strconv.Atoi(strings.TrimRight(logFile.Name(), ".log"))
			log, err := ioutil.ReadFile(logPath + "/" + logFile.Name())
			if err != nil {
				return lastStep, err
			}
			info := fmt.Sprintf("%d<->%s", lastStep, string(log))
			err = conn.WriteJSON(info)
			if err != nil {
				return lastStep, err
			}
		}
	}
	return lastStep, nil
}

func (s DeployLogService) SaveToFile(logPath string, key string, value []byte) error {
	if has, _ := common.PathExists(logPath); !has {
		if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
			return err
		}
	}
	logFile := logPath + fmt.Sprintf("%s.log", key)
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	n, _ := f.Seek(0, 2)
	_, err = f.WriteAt(value, n)
	if err != nil {
		return err
	}
	return nil
}

func (s DeployLogService) GetContainerLog(conn *websocket.Conn, containerId string) error {
	logs, err := s.containerService.GetContainerLog(containerId)
	if err != nil {
		return err
	}
	r := bufio.NewReader(logs)
	for {
		//循环从reader中根据换行符读取并转换为string
		log, err := r.ReadString('\n')
		if err != nil {
			logrus.Error(err)
			break
		}
		//todo: 去除无用字节前缀
		fmt.Println(log)
		conn.WriteJSON(string([]byte(log)[8:]))
	}

	return nil
}

func (s DeployLogService) getLogPath(project *model.Project, deployId uint) string {
	basePath := fmt.Sprintf("%s/%s/%s/%s", viper.GetString("root"), project.UserName, project.Branch, project.Name)
	return basePath + model.LogfileBasePath + fmt.Sprintf("%d/", deployId)
}
