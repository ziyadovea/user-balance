package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Error - структура для ошибки
type Error struct {
	ErrorMessage string `json:"error_message"`
}

// newErrorResponse - функция-обертка для выдачи ошибки серверу
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error().Msg(message)
	c.AbortWithStatusJSON(statusCode, &Error{ErrorMessage: message})
}
