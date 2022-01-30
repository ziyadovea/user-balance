package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Просмотр баланса
func (h *Handler) getUserAccountBalance(c *gin.Context) {

	type Request struct {
		UserID int64 `json:"user_id"`
	}

	req := &Request{}
	if err := c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect request body: %s", err.Error()))
		return
	}

	_, err := h.Services.GetUserByID(req.UserID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "this user does not exist in the system")
		return
	}

	var userBalance string

	// Проверяем, в какой валюте надо передать результат
	currency := strings.ToUpper(c.Query("currency"))
	if currency == "" || currency == "RUB" { // Рубли

		currency = "RUB"

		res, err := h.Services.GetBalanceByUserID(req.UserID, 1)
		if err != nil {
			if len(res) == 0 {
				newErrorResponse(c, http.StatusNotFound, "this user does not have a bank account")
			} else {
				newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error with getting user balance: %s", err.Error()))
			}
			return
		}
		userBalance = res

	} else { // Другая валюта

		factor, err := getFactorByCurrency(currency)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest,
				fmt.Sprintf("error with converting currency %s to RUB: %s",
					currency,
					err.Error()))
		}

		res, err := h.Services.GetBalanceByUserID(req.UserID, factor)
		if err != nil {
			if len(res) == 0 {
				newErrorResponse(c, http.StatusNotFound, "this user does not have a bank account")
			} else {
				newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error with getting user balance: %s", err.Error()))
			}
		}
		userBalance = res

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": userBalance + " " + strings.ToUpper(currency),
	})

}

// Пополнение
func (h *Handler) depositUserAccount(c *gin.Context) {

	type Request struct {
		UserID int64  `json:"user_id"`
		Amount string `json:"amount"`
	}

	// req := &Request{}
}

// Снятие
func (h *Handler) withdrawUserAccount(c *gin.Context) {

	type Request struct {
		UserID int64  `json:"user_id"`
		Amount string `json:"amount"`
	}

	// req := &Request{}
}

// Перевод
func (h *Handler) transferMoneyBetweenUsers(c *gin.Context) {

	type Request struct {
		FromUserID int64  `json:"from_user_id"`
		ToUserID   int64  `json:"to_user_id"`
		Amount     string `json:"amount"`
	}

	// req := &Request{}
}

// История транзакций
func (h *Handler) getUserTransactionHistory(c *gin.Context) {

	type Request struct {
		UserID int64 `json:"user_id"`
	}

	// req := &Request{}
}
