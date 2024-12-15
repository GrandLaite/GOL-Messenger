package db

import (
	"database/sql"
	"fmt"
	"gol_messenger/config"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось установить связь с базой данных: %w", err)
	}

	fmt.Println("Подключение к базе данных прошло успешно")
	return &Database{Connection: db}, nil
}

func (db *Database) Close() error {
	return db.Connection.Close()
}
