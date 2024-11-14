package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func connectDB() (*sql.DB, error) {
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

func main() {
	db, err := connectDB()
	if err != nil {
		fmt.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer db.Close()
	
	http.HandleFunc("/", helloHandler)
	fmt.Println("Сервер запущен и доступен по адресу http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
