# __[Тестовое задание на позицию стажера-бекендера](https://github.com/avito-tech/autumn-2021-intern-assignment)__

# __Микросервис для работы с балансом пользователей__

### Используемые технологии:
* Язык разработки - __Go__
* Реляционная СУБД - __PostgreSQL__

### Основная функциональность по работе с балансом:
* __GET /bank_account/users/:user_id/balance__ - получение баланса пользователя с id, равным :user_id <br>

Пример ответа сервера:
```json
{
    "balance": "100.00 RUB"
}
```

* __POST /bank_account/users/:user_id/balance/deposit__ - пополнение баланса пользователя с id, равным :user_id <br>

Пример тела запроса:
```json
{
    "amount": 100,
    "details": "стипендия"
}
```

Пример ответа сервера:
```json
{
    "message": "100.00 RUB were successfully deposited to user bank account"
}
```

* __POST /bank_account/users/:user_id/balance/withdraw__ - снятие с баланса пользователя с id, равным :user_id <br>

Пример тела запроса:
```json
{
    "amount": 30,
    "details": "подписка на кино"
}
```

Пример ответа сервера:
```json
{
    "message": "30.00 RUB were successfully withdrawed from user bank account"
}
```

* __POST /bank_account/users/:user_id/balance/transfer__ - перевод между пользователями (отправитель - пользователь с :user_id, получатель - в теле запроса) <br>

Пример тела запроса:
```json
{
    "to_user_id": 2,
    "amount": 4999.99,
    "details": "купи себе кроссовки"
}
```

Пример ответа сервера:
```json
{
    "message": "4999.99 RUB were successfully transferred between users"
}
```

---

### Дополнительные задания:
1) Добавить к методу получения баланса доп. параметр. <br>
Строка запроса может иметь параметр, который определяет валюту - пример: __GET /bank_account/users/:user_id/balance?currency=USD__ 

Пример ответа сервера:
```json
{
    "balance": "1.29 USD"
}
```

Возможный список валют см. на https://exchangeratesapi.io/ - API, с которого берется курс валют.

2) Необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате.
* __GET /bank_account/users/:user_id/balance/transaction_history__ - выдача истории транзакций.
Возможный параметры: __?sort=amount__ - ответ сортируется по сумме или __?sort=date__ - ответ сортируется по дате

Пример ответа сервера (здесь есть минус - баланс в копейках):
```json
[
    {
        "start_balance": 0,
        "end_balance": 10000,
        "amount": 10000,
        "message": "deposit 100.00 to account: стипендия",
        "date": "2022-01-31T22:33:39.79286Z"
    },
    {
        "start_balance": 10000,
        "end_balance": 7000,
        "amount": 3000,
        "message": "withdraw 30.00 from account: подписка на кино",
        "date": "2022-01-31T22:34:24.356805Z"
    },
    {
        "start_balance": 7000,
        "end_balance": 506999,
        "amount": 499999,
        "message": "transfer between accounts: deposit 4999.99 to account: купи себе кроссовки",
        "date": "2022-01-31T22:36:14.271042Z"
    }
]
```

---

### Дополнительная функциональность для работы с пользователями:
* __POST /users/auth/sign-up__ - регистрация нового пользователя. Изначально в БД 7 пользователей, но данный эндпоинт позволят добавить еще. <br>

Пример тела запроса:
```json
{
    "name": "Test",
    "email": "test@ex.com"
}
```
Пример ответа сервера:
```json
{
    "id": 2
}
```
* __GET /users__ - просмотр всех существующих пользователей в системе. <br>

Пример ответа сервера:
```json
[   
    {
        "id": 1,
        "name": "Иванов Иван Иванович",
        "email": "IvanovII@example.com"
    },
    {
        "id": 2,
        "name": "Test",
        "email": "test@ex.com"
    }
]
``` 

---

### Примечания к проекту:
* Архитектура REST API
* Фреймворк [gin-gonic/gin](https://github.com/gin-gonic/gin)
* Чистая архитектура (handler -> service -> repository)
* Для работы с БД используется пакет [sqlx](https://github.com/jmoiron/sqlx). Генерация файлов миграций. 
* Docker, Docker-compose, Makefile 
* Конфигурация приложения с помощью пакета [viper]("https://github.com/spf13/viper"). Работа с переменными окружения.
* Graceful Shutdown
* В БД и внутри структур языка Go баланс хранится в целых числах, чтобы не терять точность при работе с вещественными числами
* Приложение запускается на порту 8080
* Тестирование проводилось через Postman

### Схема БД: <br>
__user <-> bank_account - one-to-one <br>__
__user <-> transactions_history - one-to-many <br>__
![db-schema](https://github.com/ziyadovea/user-balance/blob/main/images/db-schema.png)

---

### Для запуска приложения используется docker-compose, вызов команд описан в makefile:

```
make build && make run
```

Если приложение запускается впервые, требуется также применить команду для миграции БД:

```
make migrate-up
```

