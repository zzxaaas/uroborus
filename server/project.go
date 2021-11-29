package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/common/auth"
	"uroborus/model"
	"uroborus/service"
)

// ProjectServer 健康检查
type ProjectServer struct {
	projectService *service.ProjectService
}

// NewProjectServer -
func NewProjectServer(projectService *service.ProjectService) *ProjectServer {
	return &ProjectServer{
		projectService: projectService,
	}
}

func (s ProjectServer) Register(c *gin.Context) {
	req := model.Project{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.projectService.Save(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s ProjectServer) CheckOut(c *gin.Context) {
	req := model.Project{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.projectService.CheckOut(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
