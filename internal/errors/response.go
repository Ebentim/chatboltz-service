package errors

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleError(c *gin.Context, err error, context string) {
	LogError(err, context)

	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, gin.H{
			"error": appErr.Message,
			"type":  appErr.Type,
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"error": "Resource not found",
			"type":  NotFoundError,
		})
		return
	}

	c.JSON(500, gin.H{
		"error": "Internal server error",
		"type":  InternalError,
	})
}
