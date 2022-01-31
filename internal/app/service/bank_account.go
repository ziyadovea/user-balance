package service

import (
	"errors"
	"github.com/ziyadovea/user-balance/internal/app/repository"
	"github.com/ziyadovea/user-balance/model"
)

// BankAccountService - структура, методы которой реализуют логику для работы с банковским счетом
type BankAccountService struct {
	repo repository.BankAccount
}

// NewBankAccountService - конструктор для BankAccountService
func NewBankAccountService(repo repository.BankAccount) *BankAccountService {
	return &BankAccountService{repo: repo}
}

// GetBalanceByUserID возвращает баланс пользователя с ID, равным userID
func (b *BankAccountService) GetBalanceByUserID(userID int64, factor float64) (string, error) {
	return b.repo.GetBalanceByUserID(userID, factor)
}

// DepositMoneyToUser начисляет amount денег пользователяю с userID
func (b *BankAccountService) DepositMoneyToUser(userID int64, amount float64, details string) error {
	if amount <= 0 {
		return errors.New("the amount must be greater than zero")
	}
	return b.repo.DepositMoneyToUser(userID, amount, details)
}

// WithdrawMoneyFromUser снимает amount денег пользователя с userID
func (b *BankAccountService) WithdrawMoneyFromUser(userID int64, amount float64, details string) error {
	if amount <= 0 {
		return errors.New("the amount must be greater than zero")
	}
	return b.repo.WithdrawMoneyFromUser(userID, amount, details)
}

// GetTransactionsHistory возвращает историю транзакций пользователя с userID
func (b *BankAccountService) GetTransactionsHistory(userID int64) ([]*model.TransactionsHistory, error) {
	return b.repo.GetTransactionsHistory(userID)
}
