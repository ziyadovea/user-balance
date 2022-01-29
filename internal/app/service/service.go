package service

import "github.com/ziyadovea/user-balance/model"

// User - интерфейс для списка методов с пользователями в слое сервиса
type User interface {
	CreateUser(*model.User) (int64, error)
	GetAllUsers() ([]*model.User, error)
}

// BankAccount - интерфейс для списка методов с банковским счетом пользователя в слое сервиса
type BankAccount interface {
}

// Service — отвечает за бизнес логику и ее переиспользование между компонентами
type Service struct {
	User
	BankAccount
}
