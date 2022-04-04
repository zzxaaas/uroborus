package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/common/auth"
	"uroborus/model"
	"uroborus/service"
)

// GroupServer -
type ProjectCommentServer struct {
	commentService *service.ProjectCommentService
}

// NewGroupServer -
func NewProjectCommentServer(projectService *service.ProjectCommentService) *ProjectCommentServer {
	return &ProjectCommentServer{
		commentService: projectService,
	}
}

func (s ProjectCommentServer) Register(c *gin.Context) {
	req := model.ProjectComment{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.FromUser = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.commentService.Register(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, req)
}

func (s ProjectCommentServer) Find(c *gin.Context) {
	req := model.ProjectComment{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := s.commentService.Find(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
