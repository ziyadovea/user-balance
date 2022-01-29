package repository

import "github.com/ziyadovea/user-balance/model"

// User - интерфейс для списка методов с пользователями в слое репозитория
type User interface {
	CreateUser(*model.User) (int64, error)
	GetAllUsers() ([]*model.User, error)
}

// BankAccount - интерфейс для списка методов с банковским счетом пользователя в слое репозитория
type BankAccount interface {
}

// Repository — отвечает за получение данных из внешних источников, такие как база данных, api, локальное хранилище и пр.
type Repository struct {
	User
	BankAccount
}
