package service

import (
	"github.com/ziyadovea/user-balance/internal/app/repository"
	"github.com/ziyadovea/user-balance/model"
)

// UserService - структура, методы которой реализуют логику для работы с пользователями
type UserService struct {
	repo repository.User
}

// NewUserService - конструктор для UserService
func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

// CreateUser создает пользователя и возвращает либо его id, либо 0 и ошибку
func (us *UserService) CreateUser(user *model.User) (int64, error) {
	if err := user.Validate(); err != nil {
		return 0, err
	}
	return us.repo.CreateUser(user)
}

// GetAllUsers позволяет посмотреть всех существующих в системе пользователей
// Возвращает либо список пользователей, либо nil и ошибку
func (us *UserService) GetAllUsers() ([]*model.User, error) {
	return us.repo.GetAllUsers()
}
