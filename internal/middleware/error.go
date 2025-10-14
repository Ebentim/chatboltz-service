package middleware

import (
	"fmt"
	"log"
	"net/http"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("[PANIC] %s: %s", c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"type":  appErrors.InternalError,
			})
		}
		c.Abort()
	})
}

func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	})
}
