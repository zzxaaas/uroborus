package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/common/auth"
	"uroborus/model"
	"uroborus/service"
)

// GroupServer -
type GroupServer struct {
	groupService *service.GroupService
}

// NewGroupServer -
func NewGroupServer(groupService *service.GroupService) *GroupServer {
	return &GroupServer{
		groupService: groupService,
	}
}

func (s GroupServer) Register(c *gin.Context) {
	req := model.Group{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.CreateUser = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.groupService.Register(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s GroupServer) Find(c *gin.Context) {
	req := model.Group{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.CreateUser = c.GetString(auth.IDTokenSubjectContextKey)
	resp, err := s.groupService.Find(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s GroupServer) Delete(c *gin.Context) {
	req := model.Group{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	req.CreateUser = c.GetString(auth.IDTokenSubjectContextKey)
	if err := s.groupService.Delete(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
