package server

import "uroborus/service"

// ProjectServer 健康检查
type ContainerServer struct {
	containerService *service.ContainerService
}

// NewProjectServer -
func NewContainerServer(containerService *service.ContainerService) *ContainerServer {
	return &ContainerServer{
		containerService: containerService,
	}
}
