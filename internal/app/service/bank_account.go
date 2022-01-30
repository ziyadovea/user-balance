package service

import "github.com/ziyadovea/user-balance/internal/app/repository"

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
