package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ziyadovea/user-balance/model"
	"strings"
)

// UserPostgres - структура, которая отвечает за свзяь работы с пользователями и PostgreSQL
type UserPostgres struct {
	db *sqlx.DB
}

// CreateUser создает пользователя и возвращает либо его id, либо 0 и ошибку
func (up *UserPostgres) CreateUser(user *model.User) (int64, error) {
	err := up.db.QueryRow(
		`INSERT INTO USERS (name, email) VALUES ($1, $2) RETURNING ID`,
		user.Name,
		user.Email,
	).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "ER_DUP_ENTRY") {
			return 0, fmt.Errorf("user with email %s already exists in the system", user.Email)
		} else {
			return 0, err
		}
	}

	return user.ID, nil
}

// GetAllUsers позволяет посмотреть всех существующих в системе пользователей
// Возвращает либо список пользователей, либо nil и ошибку
func (up *UserPostgres) GetAllUsers() ([]*model.User, error) {
	rows, err := up.db.Queryx("SELECT (id, name, email) FROM USERS")
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0)
	u := &model.User{}
	for rows.Next() {
		if err = rows.StructScan(u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
