package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ziyadovea/user-balance/model"
	"net/http"
)

// Регистрация нового пользователя
func (h *Handler) signUp(c *gin.Context) {
	user := &model.User{}
	if err := c.BindJSON(user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect request body: %s", err.Error()))
		return
	}

	id, err := h.Services.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error creating user: %s", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// Просмотр всех существующих пользователей
func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.Services.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error getting users: %s", err.Error()))
	}

	c.JSON(http.StatusOK, users)
}
