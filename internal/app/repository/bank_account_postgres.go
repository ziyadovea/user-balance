package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ziyadovea/user-balance/model"
)

// BankAccountPostgres - структура, которая отвечает за связь работы с балансом в PostgreSQL
type BankAccountPostgres struct {
	db *sqlx.DB
}

// NewBankAccountPostgres - конструктор для BankAccountPostgres
func NewBankAccountPostgres(db *sqlx.DB) *BankAccountPostgres {
	return &BankAccountPostgres{db: db}
}

// GetBalanceByUserID возвращает баланс пользователя с ID, равным userID
func (b *BankAccountPostgres) GetBalanceByUserID(userID int64, factor float64) (string, error) {
	bankAccount := &model.BankAccount{}
	err := b.db.Get(&bankAccount, "SELECT * FROM bank_account WHERE user_id=$1", userID)
	if err != nil {
		return "", err
	}

	userBalance := float64(bankAccount.Balance) / float64(100) * factor
	return fmt.Sprintf("%.2f", userBalance), nil
}
