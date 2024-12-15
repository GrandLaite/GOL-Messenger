package main

import (
	"fmt"
	"gol_messenger/auth"
	"gol_messenger/config"
	"gol_messenger/db"
	"gol_messenger/messages"
	"gol_messenger/routes"
	"gol_messenger/users"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		return
	}

	database, err := db.NewDatabase(cfg)
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("Ошибка закрытия подключения к базе данных: %v\n", err)
		}
	}()

	tokenService := auth.NewTokenService("9a8cfe1d5f3a4b2b7e6d9c2a0f1e4g3h")
	userRepository := users.NewUserRepository(database.Connection)
	userService := users.NewUserService(userRepository, tokenService)
	messageRepository := messages.NewMessageRepository(database.Connection)
	messageService := messages.NewMessageService(messageRepository)

	userHandler := users.NewUserHandler(userService)
	messageHandler := messages.NewMessageHandler(messageService)

	authMiddleware := auth.NewAuthMiddleware(tokenService)

	router := routes.SetupRoutes(userHandler, messageHandler, authMiddleware)

	port := cfg.ServerPort
	fmt.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
