package main

import (
	"fmt"
	"gol_messenger/config"
	"gol_messenger/database"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("db_config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := database.NewDatabase(cfg.DBConnectionString)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}

	fmt.Println("Сервер запущен и работает по адресу http://localhost:8080")
}
