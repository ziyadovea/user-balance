package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Просмотр баланса
func (h *Handler) getUserAccountBalance(c *gin.Context) {

	// Достаем id пользователя из url
	userID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	// Проверяем существование пользователя в системе
	_, err = h.Services.GetUserByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "this user does not exist in the system")
		return
	}

	var userBalance string

	// Проверяем, в какой валюте надо передать результат
	currency := strings.ToUpper(c.Query("currency"))
	if currency == "" || currency == "RUB" { // Рубли

		currency = "RUB"

		res, err := h.Services.GetBalanceByUserID(userID, 1)
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
				fmt.Sprintf("error with converting RUB to currency %s: %s",
					currency,
					err.Error()))
			return
		}

		res, err := h.Services.GetBalanceByUserID(userID, factor)
		if err != nil {
			if len(res) == 0 {
				newErrorResponse(c, http.StatusNotFound, "this user does not have a bank account")
			} else {
				newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error with getting user balance: %s", err.Error()))
			}
			return
		}
		userBalance = res

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": userBalance + " " + strings.ToUpper(currency),
	})

}

// Пополнение или снятие
func (h *Handler) updateUserAccount(c *gin.Context) {

	// Достаем id пользователя из url
	userID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	// Тело запроса
	type Request struct {
		Amount  float64 `json:"amount"`
		Details string  `json:"details"`
	}

	// Декодируем тело запроса
	req := &Request{}
	if err = c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect request body: %s", err.Error()))
		return
	}

	// Проверяем существования пользователя в системе
	_, err = h.Services.GetUserByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "this user does not exist in the system")
		return
	}

	// В зависимости от адреса - пополненение или снятие денег
	isDeposit := false                                     // Флаг для вывода ответа
	if strings.Contains(c.Request.RequestURI, "deposit") { // Пополнене
		isDeposit = true
		err = h.Services.DepositMoneyToUser(userID, req.Amount, req.Details)
	} else if strings.Contains(c.Request.RequestURI, "withdraw") { // Снятие
		err = h.Services.WithdrawMoneyFromUser(userID, req.Amount, req.Details)
	}
	if err != nil {
		if strings.Contains(err.Error(), "the amount must be greater than zero") {
			newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error with update user balance: %s", err.Error()))
		} else {
			newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error with update user balance: %s", err.Error()))
		}
		return
	}

	// В зависимости от операции - разный ответ
	if isDeposit {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("%.2f RUB were successfully deposited to user bank account", req.Amount),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("%.2f RUB were successfully withdrawed from user bank account", req.Amount),
		})
	}
}

// Перевод
func (h *Handler) transferMoneyBetweenUsers(c *gin.Context) {

	// Достаем id отправителя из url
	fromUserID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	// Тело запроса
	type Request struct {
		ToUserID int64   `json:"to_user_id"`
		Amount   float64 `json:"amount"`
		Details  string  `json:"details"`
	}

	// Декодируем тело запроса
	req := &Request{}
	if err := c.BindJSON(req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("incorrect request body: %s", err.Error()))
		return
	}

	// Проверяем существование юзера-отправителя
	_, err = h.Services.GetUserByID(fromUserID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "sender user does not exist in the system")
		return
	}

	// Проверяем существования юзера-получателя
	_, err = h.Services.GetUserByID(req.ToUserID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "destination user does not exist in the system")
		return
	}

	// Делегируем работу слою сервиса
	err = h.Services.TransferMoneyBetweenUsers(fromUserID, req.ToUserID, req.Amount, req.Details)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error with transfer money between users: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%.2f RUB were successfully transferred between users", req.Amount),
	})
}

// История транзакций
func (h *Handler) getUserTransactionHistory(c *gin.Context) {

	// Достаем id пользователя из url
	userID, err := h.getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	// Проверяем существования пользователя в системе
	_, err = h.Services.GetUserByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "this user does not exist in the system")
		return
	}

	// Получаем историю транзакций
	history, err := h.Services.GetTransactionsHistory(userID)
	if len(history) == 0 {
		newErrorResponse(c, http.StatusNotFound, "there is no transaction history for this user")
		return
	}

	// Сортируем результат
	sortParam := c.Query("sort")
	switch strings.ToLower(sortParam) {
	case "amount":
		sort.Slice(history, func(i, j int) bool {
			return history[i].Amount < history[j].Amount
		})
	case "date":
		sort.Slice(history, func(i, j int) bool {
			return history[i].Date.Unix() < history[j].Date.Unix()
		})
	case "":
	default:
		newErrorResponse(c, http.StatusBadRequest, "invalid sort param in query")
		return
	}

	c.JSON(http.StatusOK, history)
}

// Достаем id пользователя из url
func (h *Handler) getUserID(c *gin.Context) (int64, error) {
	strUserID := c.Param("user_id")
	intUserID, err := strconv.Atoi(strUserID)
	if err != nil {
		return 0, err
	}
	return int64(intUserID), nil
}
