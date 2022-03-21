package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"uroborus/common/auth"
	"uroborus/model"
	"uroborus/service"
)

// ProjectServer 健康检查
type ContainerServer struct {
	containerService *service.ContainerService
	upgrader         websocket.Upgrader
}

// NewProjectServer -
func NewContainerServer(containerService *service.ContainerService) *ContainerServer {
	return &ContainerServer{
		containerService: containerService,
		upgrader: websocket.Upgrader{
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s ContainerServer) GetAll(c *gin.Context) {
	user := c.GetString(auth.IDTokenSubjectContextKey)
	err, resp := s.containerService.GetAll(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s ContainerServer) Exec(c *gin.Context) {
	req := model.ConnectContainerReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := s.containerService.Terminal(ws, req); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
}
