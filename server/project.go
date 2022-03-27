package server

import (
	"errors"
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

func (s ProjectServer) Get(c *gin.Context) {
	req := model.Project{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if resp, err := s.projectService.Find(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func (s ProjectServer) Register(c *gin.Context) {
	req := model.RegisterProjectReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.projectService.Save(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": req.ID})
}

func (s ProjectServer) Delete(c *gin.Context) {
	req := model.Project{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.projectService.Delete(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s ProjectServer) GetGroupProjs(c *gin.Context) {
	req := model.Group{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := s.projectService.GetGroupProjects(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s ProjectServer) RegisterGroupProj(c *gin.Context) {
	req := model.RegisterGroupProjectReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user := c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.projectService.RegisterGroupProject(&req, user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s ProjectServer) SaveProjectInfo(c *gin.Context) {
	f, err := c.FormFile("img")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req := model.Project{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.UserName = c.GetString(auth.IDTokenSubjectContextKey)
	filePath, err := s.projectService.SaveProjectInfo(&req, f)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if f != nil {
		c.SaveUploadedFile(f, filePath)
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
