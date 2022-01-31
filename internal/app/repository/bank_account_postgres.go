package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ziyadovea/user-balance/model"
	"math"
	"time"
)

// BankAccountPostgres - структура, которая отвечает за связь работы с балансом в PostgreSQL
type BankAccountPostgres struct {
	db *sqlx.DB
}

// NewBankAccountPostgres - конструктор для BankAccountPostgres
func NewBankAccountPostgres(db *sqlx.DB) *BankAccountPostgres {
	return &BankAccountPostgres{db: db}
}

// GetBalanceByUserID возвращает баланс пользователя с ID, равным userID
func (b *BankAccountPostgres) GetBalanceByUserID(userID int64, factor float64) (string, error) {
	bankAccount := &model.BankAccount{}
	err := b.db.Get(bankAccount, "SELECT * FROM bank_account WHERE user_id=$1", userID)
	if err != nil {
		return "", err
	}

	userBalance := float64(bankAccount.Balance) / float64(100) * factor
	return fmt.Sprintf("%.2f", userBalance), nil
}

// DepositMoneyToUser начисляет amount денег пользователяю с userID
func (b *BankAccountPostgres) DepositMoneyToUser(userID int64, amount float64, details string) error {

	// Сразу переведем деньги в целые числа (храним в копейках)
	intAmount := int64(math.Round(amount * 100))

	// Сразу сформируем сообщение для записи в таблицу транзакций
	msg := fmt.Sprintf("deposit %.2f to account", amount)
	if len(details) != 0 {
		msg += ": " + details
	}

	bankAccount := &model.BankAccount{}
	err := b.db.Get(bankAccount, "SELECT * FROM bank_account WHERE user_id=$1", userID)

	// Кейс 1: у пользователя еще нет записи баланса в таблице
	// значит надо создать запись в БД о балансе такого пользователя
	if err != nil {
		return b.firstDeposit(userID, intAmount, msg)
	}

	// Кейс 2: у пользователя уже есть запись в таблице
	// значит надо обновить это значение

	return b.notFirstDeposit(userID, intAmount, msg, bankAccount.Balance)
}

// WithdrawMoneyFromUser снимает amount денег пользователяю с userID
func (b *BankAccountPostgres) WithdrawMoneyFromUser(userID int64, amount float64, details string) error {

	// Сразу переведем деньги в целые числа (храним в копейках)
	intAmount := int64(math.Round(amount * 100))

	bankAccount := &model.BankAccount{}
	err := b.db.Get(bankAccount, "SELECT * FROM bank_account WHERE user_id=$1", userID)

	// Кейс 1: у пользователя еще нет записи баланса в таблице
	// Ошибка, так как снимать не с чего
	if err != nil {
		return errors.New("this user does not have a bank account")
	}

	// Кейс 2: у пользователя уже есть запись в таблице
	// значит надо обновить это значение

	// Проверим, достаточно ли у пользователя денег
	if intAmount > bankAccount.Balance {
		return errors.New("not enough money in the bank account")
	}

	// Сформируем сообщение для записи в таблицу транзакций
	msg := fmt.Sprintf("withdraw %.2f from account", amount)
	if len(details) != 0 {
		msg += ": " + details
	}

	return b.notFirstDeposit(userID, -intAmount, msg, bankAccount.Balance)
}

// TransferMoneyBetweenUsers переводит amount денег со счета пользователя с fromUserID пользователю с toUserID
func (b *BankAccountPostgres) TransferMoneyBetweenUsers(fromUserID int64, toUserID int64, amount float64, details string) error {

	// Достаем из БД запись о банковском счете пользователя-отправителя
	fromBankAccount := &model.BankAccount{}
	err := b.db.Get(fromBankAccount, "SELECT * FROM bank_account WHERE user_id=$1", fromUserID)

	// У пользователя еще нет записи баланса в таблице
	// Ошибка, так как снимать не с чего
	if err != nil {
		return errors.New("sender user does not have a bank account")
	}

	// Переведем деньги в целые числа (храним в копейках)
	intAmount := int64(math.Round(amount * 100))

	// Проверяем, достаточно ли денег у отправителя на счету
	if fromBankAccount.Balance < intAmount {
		return errors.New("the sender does not have enough money in the account")
	}

	// Сформируем сообщение для записи в таблицу транзакций для отправителя
	fromMsg := fmt.Sprintf("transfer between accounts: withdraw %.2f from account", amount)
	if len(details) != 0 {
		fromMsg += ": " + details
	}

	// Сформируем сообщение для записи в таблицу транзакций для получателя
	toMsg := fmt.Sprintf("transfer between accounts: deposit %.2f to account", amount)
	if len(details) != 0 {
		toMsg += ": " + details
	}

	// Достаем из БД запись о банковском счете пользователя-получателя
	toBankAccount := &model.BankAccount{}
	err = b.db.Get(toBankAccount, "SELECT * FROM bank_account WHERE user_id=$1", toUserID)
	// Кейс 1
	// У пользователя еще нет записи баланса в таблице
	// Значит надо снять у отправителя, а получателю создать запись
	if err != nil {
		return b.firstTransfer(fromUserID, toUserID, fromMsg, toMsg, fromBankAccount.Balance, intAmount)
	}

	// Кейс 2
	// У пользователя уже есть запись баланса в таблице
	// Значит надо снять у отправителя, а получателю увеличить баланс
	return b.notFirstTransfer(fromUserID, toUserID, fromMsg, toMsg, fromBankAccount.Balance, toBankAccount.Balance, intAmount)

}

// GetTransactionsHistory возвращает историю транзакций пользователя с userID
func (b *BankAccountPostgres) GetTransactionsHistory(userID int64) ([]*model.TransactionsHistory, error) {
	history := make([]*model.TransactionsHistory, 0)
	err := b.db.Select(&history, "SELECT * FROM transactions_history WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// firstDeposit реализует первый взнос денег на счет пользователя (когда у пользователя нет еще записи в БД)
func (b *BankAccountPostgres) firstDeposit(userID int64, amount int64, msg string) error {

	// Надо сделать 2 вещи: обновить баланс пользователя
	// И добавить соответствующую запись в таблицу историй транзакций
	// Объединим эти операции в транзакцию
	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO bank_account (user_id, balance) VALUES ($1, $2)",
		userID,
		amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		userID, 0, amount, amount, msg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// notFirstDeposit реализует не первый взнос денег на счет пользователя
func (b *BankAccountPostgres) notFirstDeposit(userID int64, amount int64, msg string, startBalance int64) error {
	// Надо сделать 2 вещи: обновить баланс пользователя
	// И добавить соответствующую запись в таблицу историй транзакций
	// Объединим эти операции в транзакцию
	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}

	endBalance := startBalance + amount
	_, err = tx.Exec(
		"UPDATE bank_account SET balance=$1 WHERE user_id=$2",
		endBalance,
		userID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Отрицательное значение передается при снятии денег
	// Так как в таблице есть ограничение на то, что сумма только положительная
	// Умножим ее на минус 1
	if amount < 0 {
		amount *= -1
	}
	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		userID, startBalance, endBalance, amount, msg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// firstTransfer реализует первый перевод между пользователями (когда у получателя нет еще записи в БД)
func (b *BankAccountPostgres) firstTransfer(
	fromUserID int64, toUserID int64,
	fromMsg string, toMsg string,
	fromStartBalance int64, amount int64) error {

	// Надо сделать 4 вещи:
	// 1. Снять со счета отправителя
	// 2. Добавить соответствующую запись в таблицу историй транзакций
	// 3. Создать запись в таблице получателю
	// 4. Добавить соответствующую запись в таблицу историй транзакций
	// Объединим эти операции в транзакции

	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}

	// 1. Снять со счета отправителя
	endBalance := fromStartBalance - amount
	_, err = tx.Exec(
		"UPDATE bank_account SET balance=$1 WHERE user_id=$2",
		endBalance,
		fromUserID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Добавить соответствующую запись в таблицу историй транзакций
	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		fromUserID, fromStartBalance, endBalance, amount, fromMsg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Создать запись в таблице получателю
	_, err = tx.Exec(
		"INSERT INTO bank_account (user_id, balance) VALUES ($1, $2)",
		toUserID,
		amount,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4. Добавить соответствующую запись в таблицу историй транзакций
	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		toUserID, 0, amount, amount, toMsg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// notFirstTransfer реализует не первый перевод между пользователями
func (b *BankAccountPostgres) notFirstTransfer(
	fromUserID int64, toUserID int64,
	fromMsg string, toMsg string,
	fromStartBalance int64, toStartBalance int64,
	amount int64) error {

	// Надо сделать 4 вещи:
	// 1. Снять со счета отправителя
	// 2. Добавить соответствующую запись в таблицу историй транзакций
	// 3. Положить на счет получателю
	// 4. Добавить соответствующую запись в таблицу историй транзакций
	// Объединим эти операции в транзакции

	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}

	// 1. Снять со счета отправителя
	endBalance := fromStartBalance - amount
	_, err = tx.Exec(
		"UPDATE bank_account SET balance=$1 WHERE user_id=$2",
		endBalance,
		fromUserID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Добавить соответствующую запись в таблицу историй транзакций
	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		fromUserID, fromStartBalance, endBalance, amount, fromMsg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Положить на счет получателю
	endBalance = toStartBalance + amount
	_, err = tx.Exec(
		"UPDATE bank_account set balance=$1 WHERE user_id=$2",
		endBalance,
		toUserID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4. Добавить соответствующую запись в таблицу историй транзакций
	_, err = tx.Exec(
		"INSERT INTO transactions_history (user_id, start_balance, end_balance, amount, message, date) "+
			"VALUES ($1, $2, $3, $4, $5, $6)",
		toUserID, toStartBalance, endBalance, amount, toMsg, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}
