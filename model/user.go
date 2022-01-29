package model

import "github.com/go-playground/validator/v10"

// User - структура для представления пользователя
type User struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" validate:"required" db:"name"`
	Email string `json:"email" validate:"required,email" db:"email"`
}

var validate *validator.Validate

// Validate проверяет валидность полей пользователя
func (u *User) Validate() error {
	validate = validator.New()
	return validate.Struct(u)
}
