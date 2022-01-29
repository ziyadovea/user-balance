package model

import "time"

// TransactionsHistory - структура для представления истории денежных переводов
type TransactionsHistory struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	StartBalance int64     `json:"start_balance"`
	EndBalance   int64     `json:"end_balance"`
	Amount       int64     `json:"amount"`
	Message      string    `json:"message"`
	Date         time.Time `json:"date"`
}
