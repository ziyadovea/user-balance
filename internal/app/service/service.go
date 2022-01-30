package service

import (
	"github.com/ziyadovea/user-balance/internal/app/repository"
	"github.com/ziyadovea/user-balance/model"
)

// User - интерфейс для списка методов с пользователями в слое сервиса
type User interface {
	CreateUser(user *model.User) (int64, error)
	GetAllUsers() ([]*model.User, error)
	GetUserByID(userID int64) (*model.User, error)
}

// BankAccount - интерфейс для списка методов с банковским счетом пользователя в слое сервиса
type BankAccount interface {
	GetBalanceByUserID(userID int64, factor float64) (string, error)
}

// Service — отвечает за бизнес логику и ее переиспользование между компонентами
type Service struct {
	User
	BankAccount
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:        NewUserService(repo.User),
		BankAccount: NewBankAccountService(repo),
	}
}
