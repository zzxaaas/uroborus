package middleware

import (
	"encoding/json"
	"strings"
	"uroborus/common/logging"

	"github.com/gin-gonic/gin"
)

// BodyError -
type BodyError struct {
	Message string `json:"message"`
}

// Error returns an error handler
func Error(logger *logging.ZapLogger) gin.HandlerFunc {
	// TODO: error response
	return func(c *gin.Context) {
		c.Next()
		if messages := c.Errors.Errors(); len(messages) > 0 && c.Writer.Size() == 0 {
			bodyError := BodyError{Message: strings.Join(messages, "\n")}
			if body, err := json.Marshal(bodyError); err != nil {
				logger.Error(err.Error())
			} else {
				c.Writer.Write(body)
			}
		}
	}
}
