package reusable

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsPremium bool   `json:"is_premium"`
}

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

	query := `INSERT INTO Users (username, email, password, is_premium) VALUES ($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, user.Username, user.Email, user.Password, user.IsPremium).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Ошибка регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Пользователь успешно зарегистрирован с ID: %d", user.ID)))
}

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
	query := `SELECT id, username, email, is_premium FROM Users WHERE id = $1`
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.IsPremium)
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

	query := `UPDATE Users SET username = $1, email = $2, is_premium = $3 WHERE id = $4`
	_, err = db.Exec(query, updatedUser.Username, updatedUser.Email, updatedUser.IsPremium, userID)
	if err != nil {
		http.Error(w, "Не удалось обновить данные пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Данные пользователя успешно обновлены"))
}

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
