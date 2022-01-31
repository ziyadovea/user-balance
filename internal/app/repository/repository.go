package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/ziyadovea/user-balance/model"
)

// User - интерфейс для списка методов с пользователями в слое репозитория
type User interface {
	CreateUser(user *model.User) (int64, error)
	GetAllUsers() ([]*model.User, error)
	GetUserByID(userID int64) (*model.User, error)
}

// BankAccount - интерфейс для списка методов с банковским счетом пользователя в слое репозитория
type BankAccount interface {
	GetBalanceByUserID(userID int64, factor float64) (string, error)
	DepositMoneyToUser(userID int64, amount float64, details string) error
	WithdrawMoneyFromUser(userID int64, amount float64, details string) error
	TransferMoneyBetweenUsers(fromUserID int64, toUserID int64, amount float64, details string) error
	GetTransactionsHistory(userID int64) ([]*model.TransactionsHistory, error)
}

// Repository — отвечает за получение данных из внешних источников, такие как база данных, api, локальное хранилище и пр.
type Repository struct {
	User
	BankAccount
}

// NewRepository - конструктор для Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserPostgres(db),
		BankAccount: NewBankAccountPostgres(db),
	}
}
