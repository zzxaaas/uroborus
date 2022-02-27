package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"uroborus/model"
	"uroborus/service"
)

// ProjectServer 健康检查
type DeployServer struct {
	deployService        *service.DeployService
	deployHistoryService *service.DeployHistoryService
	deployLogService     *service.DeployLogService
	upgrader             websocket.Upgrader
}

// NewProjectServer -
func NewDeployServer(deployService *service.DeployService, deployHistoryService *service.DeployHistoryService, deployLogService *service.DeployLogService) *DeployServer {
	return &DeployServer{
		deployService:        deployService,
		deployHistoryService: deployHistoryService,
		deployLogService:     deployLogService,
		upgrader: websocket.Upgrader{
			// 解决跨域问题
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s DeployServer) Get(c *gin.Context) {
	req := &model.DeployHistory{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if resp, err := s.deployHistoryService.Find(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s DeployServer) Deploy(c *gin.Context) {
	req := model.DeployHistory{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := s.deployService.Deploy(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s DeployServer) Log(c *gin.Context) {
	req := model.DeployHistory{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer ws.Close()
	if err := s.deployLogService.GetLog(ws, &req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
