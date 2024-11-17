package database

import (
	"database/sql"
	"fmt"
	"gol_messenger/config"

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

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось установить связь с базой данных: %w", err)
	}

	return &Database{Connection: db}, nil
}

func (db *Database) Close() error {
	return db.Connection.Close()
}
