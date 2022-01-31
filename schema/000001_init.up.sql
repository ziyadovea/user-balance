CREATE TABLE users
(
    id    SERIAL  NOT NULL PRIMARY KEY,
    name  VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE
);

INSERT INTO users (name, email)
VALUES
    ('Иванов Иван Иванович',     'IvanovII@example.com'),
    ('Петров Петр Петрович',     'PetrovPP@example.com'),
    ('Крид Егор Блэкстарович',   'KreedEB@bs.com'),
    ('Солнцева Диана Андреевна', 'SolntsevaDA@example.com'),
    ('Шатилова Вера Кирилловна', 'ShatilovaVK@example.com'),
    ('Попугаев Егор Алексеевич', 'PopugaevEA@example.com'),
    ('Зиядова Амелия Эмильевна', 'ZiyadovaAE@example.com');

CREATE TABLE bank_account
(
    id      SERIAL NOT NULL PRIMARY KEY,
    user_id INT    NOT NULL UNIQUE,
    balance BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT balance_check CHECK (balance >= 0)
);

CREATE TABLE transactions_history
(
    id            SERIAL    NOT NULL PRIMARY KEY,
    user_id       INT       NOT NULL,
    start_balance BIGINT    NOT NULL,
    end_balance   BIGINT    NOT NULL,
    amount        BIGINT    NOT NULL,
    message       VARCHAR   NULL,
    date          TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT balance_check CHECK (start_balance >= 0 AND
                                    end_balance >= 0 AND
                                    amount >= 0)
);















