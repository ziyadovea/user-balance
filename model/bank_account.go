package model

// BankAccount - структура для представления банковского аккаунта пользователя
// Пока здесь лишь одно смысловое поле - баланс пользователя, однако при расширении
// приложения список полей может расширяться
type BankAccount struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	Balance int64 `json:"balance"`
}
