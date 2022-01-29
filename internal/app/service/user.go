package service

import (
	"github.com/ziyadovea/user-balance/internal/app/repository"
	"github.com/ziyadovea/user-balance/model"
)

// UserService - структура, методы которой реализуют логику для работы с пользователями
type UserService struct {
	repo repository.User
}

// CreateUser создает пользователя и возвращает либо его id, либо 0 и ошибку
func (us *UserService) CreateUser(user *model.User) (int64, error) {
	return us.repo.CreateUser(user)
}

// GetAllUsers позволяет посмотреть всех существующих в системе пользователей
// Возвращает либо список пользователей, либо nil и ошибку
func (us *UserService) GetAllUsers() ([]*model.User, error) {
	return us.repo.GetAllUsers()
}
