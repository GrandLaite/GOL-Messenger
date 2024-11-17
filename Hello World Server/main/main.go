package main

import (
	"fmt"
	"gol_messenger/config"
	"gol_messenger/database"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("db_config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	fmt.Println("База данных и сервер успешно запущены.")
}
