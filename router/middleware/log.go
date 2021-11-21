package middleware

import (
	"time"
	"uroborus/common/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Log middleware
func Log(logger *logging.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
			zap.Duration("latency", latency),
		}
		if err := c.Errors.Last(); err != nil {
			fields = append(fields, zap.Error(err))
		}

		logger.Info(path, fields...)
	}
}
