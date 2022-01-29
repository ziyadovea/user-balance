package handler

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	ErrorMessage string `json:"error_message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, &Error{ErrorMessage: message})
}
