package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/model"
	"uroborus/service"
)

type UserServer struct {
	userService *service.UserService
}

func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{
		userService: userService,
	}
}

func (s UserServer) Register(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := s.userService.Register(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s UserServer) Login(c *gin.Context) {
	user := model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token, err := s.userService.Login(&user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id_token": token})
}
