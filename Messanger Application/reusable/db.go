package reusable

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB() (*sql.DB, error) {
	dsn := "postgres://postgres:1@localhost:5432/gol_messenger"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось авторизовать соединение с базой данных: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	fmt.Println("Подключение к базе данных прошло успешно")
	return db, nil
}
