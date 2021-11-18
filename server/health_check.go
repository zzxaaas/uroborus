package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

// HealthServer 健康检查
type HealthServer struct {
}

// NewHealthServer constructor
func NewHealthServer() *HealthServer {
	return &HealthServer{}
}

// CheckV1 health
func (s *HealthServer) CheckV1(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service_name": "uroboros",
		"overall":      "success",
		"env":          viper.GetString("env"),
		"commit_id":    "",
		"dependencies": gin.H{},
	})
}
