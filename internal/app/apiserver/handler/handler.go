package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ziyadovea/user-balance/internal/app/service"
)

// Handler - структура для обработчика
type Handler struct {
	Services *service.Service
}

// NewHandler - конструктор для обработчика
func NewHandler(services *service.Service) *Handler {
	return &Handler{Services: services}
}

// InitRoutes - функция для инициализации обработчиков эндпоинтов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Основные роуты для работы с балансом
	balance := router.Group("/balance")
	{
		balance.GET("/", h.getUserAccountBalance)                        // Просмотр баланса
		balance.POST("/deposit", h.depositUserAccount)                   // Пополнение
		balance.POST("/withdraw", h.withdrawUserAccount)                 // Снятие
		balance.POST("/transfer", h.transferMoneyBetweenUsers)           // Перевод
		balance.GET("/transaction_history", h.getUserTransactionHistory) // История транзакций
	}

	// Дополнительные роуты для удобства работы с пользователями
	router.POST("/users/auth/sign-up", h.signUp) // Регистрация нового пользователя
	router.GET("/users", h.getAllUsers)          // Просмотр всех существующих пользователей

	return router
}
