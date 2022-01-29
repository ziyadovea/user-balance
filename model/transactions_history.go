package model

import "time"

// TransactionsHistory - структура для представления истории денежных переводов
type TransactionsHistory struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	StartBalance int64     `json:"start_balance" db:"start_balance"`
	EndBalance   int64     `json:"end_balance" db:"end_balance"`
	Amount       int64     `json:"amount" db:"amount"`
	Message      string    `json:"message" db:"message"`
	Date         time.Time `json:"date" db:"date"`
}
