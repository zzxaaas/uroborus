package service

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"uroborus/model"
)

type ProxyService struct {
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (s ProxyService) getProxyServer(access_url string) string {
	parseUrl, _ := url.Parse(access_url)
	return parseUrl.Hostname()
}

func (s ProxyService) GetConfFile() (*os.File, *os.File, error) {
	in, err := os.Open(viper.GetString("nginx.conf"))
	if err != nil {
		return nil, nil, fmt.Errorf("open file fail:%v", err)
	}
	out, err := os.OpenFile(viper.GetString("nginx.conf"), os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return nil, nil, fmt.Errorf("Open write file fail:%v", err)
	}
	return in, out, nil
}

func (s ProxyService) FoundLinePosi(in *os.File, condition string) ([]string, int) {
	file_lines := make([]string, 0)
	index := 0
	br := bufio.NewReader(in)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		file_lines = append(file_lines, string(line))
		if strings.Contains(string(line), condition) {
			break
		}
		index++
	}
	return file_lines, index
}

func (s ProxyService) RegisterProxy(req model.Project) error {
	in, out, err := s.GetConfFile()
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	fileLines, index := s.FoundLinePosi(in, model.Placeholder)
	proxy := fmt.Sprintf(model.ProxyConfig, s.getProxyServer(req.AccessUrl), req.BindPort)
	fileLines[index] = strings.ReplaceAll(fileLines[index], model.Placeholder, proxy)

	for _, line := range fileLines {
		if _, err = out.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("write to file fail:%v", err)
		}
	}
	return nil
}

func (s ProxyService) RemoveProxy(req model.Project) error {
	in, out, err := s.GetConfFile()
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	fileLines, index := s.FoundLinePosi(in, strconv.Itoa(req.BindPort))
	fileLines = append(fileLines[:index-4], fileLines[index+3:]...)

	for _, line := range fileLines {
		if _, err = out.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("write to file fail:%v", err)
		}
	}
	return nil
}
