package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер для PostgreSQL через pgx
)

// Функция для подключения к базе данных PostgreSQL
func ConnectDB() (*sql.DB, error) {
	// Строка подключения к базе данных
	dsn := "postgres://postgres:1@localhost:5432/gol_messenger" // Укажите свои параметры
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}
