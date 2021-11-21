package common

import (
	"github.com/gin-gonic/gin"
)

// 运行模式
const (
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
	TestMode    = gin.TestMode
)

// Mode 获取运行模式（使用 GIN_MODE 环境变量）
func Mode() string {
	return gin.Mode()
}
