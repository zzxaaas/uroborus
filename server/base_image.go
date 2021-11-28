package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"uroborus/model"
	"uroborus/service"
)

// HealthServer 健康检查
type BaseImageServer struct {
	baseImageService *service.BaseImageService
}

// NewHealthServer constructor
func NewBaseImageServer(baseImageService *service.BaseImageService) *BaseImageServer {
	return &BaseImageServer{
		baseImageService: baseImageService,
	}
}

func (s BaseImageServer) Register(c *gin.Context) {
	req := model.BaseImage{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	header := c.GetHeader("User-Agent")
	if err := s.baseImageService.Save(req, header); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (s BaseImageServer) Get(c *gin.Context) {
	req := model.BaseImage{}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := s.baseImageService.Get(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
