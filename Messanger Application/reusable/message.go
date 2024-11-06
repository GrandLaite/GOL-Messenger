package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
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

	// Вставка сообщения в базу данных
	query := `INSERT INTO Messages (user_id, content) VALUES ($1, $2) RETURNING id, created_at`
	err = db.QueryRow(query, msg.UserID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
	if err != nil {
		http.Error(w, "Не удалось создать сообщение", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Сообщение успешно создано с ID: %d", msg.ID)))
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	if messageID == "" {
		http.Error(w, "Необходимо указать ID сообщения", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var msg Message
	query := `SELECT id, user_id, content, created_at FROM Messages WHERE id = $1`
	err = db.QueryRow(query, messageID).Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Сообщение не найдено", http.StatusNotFound)
		} else {
			http.Error(w, "Не улаорст получить сообщение", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	if messageID == "" {
		http.Error(w, "Необходимо указать ID сообщения", http.StatusBadRequest)
		return
	}

	var updatedMsg Message
	err := json.NewDecoder(r.Body).Decode(&updatedMsg)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `UPDATE Messages SET content = $1 WHERE id = $2`
	_, err = db.Exec(query, updatedMsg.Content, messageID)
	if err != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message updated successfully"))
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `DELETE FROM Messages WHERE id = $1`
	_, err = db.Exec(query, messageID)
	if err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message deleted successfully"))
}

func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, user_id, content, created_at FROM Messages`)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			http.Error(w, "Error scanning messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
