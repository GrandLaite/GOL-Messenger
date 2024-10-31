package main

import (
	"fmt"
	"net/http"
)

func main() {
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	fmt.Println("Connected to the database successfully")

	http.HandleFunc("/register", RegisterUserHandler)                  // Регистрация нового пользователя
	http.HandleFunc("/user", AuthMiddleware(GetUserHandler))           // Получение информации о пользователе
	http.HandleFunc("/user/update", AuthMiddleware(UpdateUserHandler)) // Обновление данных пользователя
	http.HandleFunc("/user/delete", AuthMiddleware(DeleteUserHandler)) // Удаление пользователя

	http.HandleFunc("/messages", AuthMiddleware(ListMessagesHandler))        // Список всех сообщений
	http.HandleFunc("/message", AuthMiddleware(GetMessageHandler))           // Получение сообщения по ID
	http.HandleFunc("/message/create", AuthMiddleware(CreateMessageHandler)) // Создание сообщения
	http.HandleFunc("/message/update", AuthMiddleware(UpdateMessageHandler)) // Обновление сообщения
	http.HandleFunc("/message/delete", AuthMiddleware(DeleteMessageHandler)) // Удаление сообщения

	http.HandleFunc("/login", LoginHandler) // Аутентификация и генерация JWT

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
