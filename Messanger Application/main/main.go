package main

import (
	"fmt"
	"gol_messenger/reusable"
	"net/http"
)

func main() {
	db, err := reusable.ConnectDB()
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/register", reusable.RegisterUserHandler)                           // Регистрация нового пользователя
	http.HandleFunc("/login", reusable.LoginHandler)                                     // Аутентификация и генерация JWT
	http.HandleFunc("/user", reusable.AuthMiddleware(reusable.GetUserHandler))           // Получение информации о пользователе
	http.HandleFunc("/user/update", reusable.AuthMiddleware(reusable.UpdateUserHandler)) // Обновление данных пользователя
	http.HandleFunc("/user/delete", reusable.AuthMiddleware(reusable.DeleteUserHandler)) // Удаление пользователя

	http.HandleFunc("/messages", reusable.AuthMiddleware(reusable.ListMessagesHandler))              // Список всех сообщений
	http.HandleFunc("/message", reusable.AuthMiddleware(reusable.GetMessageHandler))                 // Получение сообщения по ID
	http.HandleFunc("/message/create", reusable.AuthMiddleware(reusable.CreateMessageHandler))       // Создание сообщения
	http.HandleFunc("/message/update", reusable.AuthMiddleware(reusable.UpdateMessageHandler))       // Обновление сообщения (только для премиум-пользователей)
	http.HandleFunc("/message/delete", reusable.AuthMiddleware(reusable.DeleteMessageHandler))       // Удаление сообщения (только для премиум-пользователей)
	http.HandleFunc("/message/like", reusable.AuthMiddleware(reusable.LikeMessageHandler))           // Лайк сообщения
	http.HandleFunc("/message/superlike", reusable.AuthMiddleware(reusable.SuperlikeMessageHandler)) // Суперлайк сообщения (только для премиум-пользователей)

	fmt.Println("Сервер запущен и работает по адресу http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
