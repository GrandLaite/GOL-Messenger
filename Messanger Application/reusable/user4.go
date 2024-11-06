package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

// Регистрация пользователя
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Создание пользователя в базе данных
	query := `INSERT INTO Users (username, email, password, role_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, user.Username, user.Email, user.Password, user.RoleID).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Пользователь успешно зарегистрирован с ID: %d", user.ID)))
}

// Получение информации о пользователе по ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Требуется ID пользователя", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user User
	query := `SELECT id, username, email, role_id FROM Users WHERE id = $1`
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Не удалось получить данные пользователя", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Обновление информации о пользователе
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Требуется ID пользователя", http.StatusBadRequest)
		return
	}

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Некорректный формат запроса", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `UPDATE Users SET username = $1, email = $2 WHERE id = $3`
	_, err = db.Exec(query, updatedUser.Username, updatedUser.Email, userID)
	if err != nil {
		http.Error(w, "Не удалось обновить данные пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные пользователя успешно обновлены"))
}

// Удаление пользователя
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "Требуется ID пользователя", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `DELETE FROM Users WHERE id = $1`
	_, err = db.Exec(query, userID)
	if err != nil {
		http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно удален"))
}
