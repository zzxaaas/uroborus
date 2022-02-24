package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/model"
	"uroborus/service"
)

// ProjectServer 健康检查
type DeployServer struct {
	deployService        *service.DeployService
	deployHistoryService *service.DeployHistoryService
}

// NewProjectServer -
func NewDeployServer(deployService *service.DeployService, deployHistoryService *service.DeployHistoryService) *DeployServer {
	return &DeployServer{
		deployService:        deployService,
		deployHistoryService: deployHistoryService,
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
	//req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.deployService.Deploy(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
