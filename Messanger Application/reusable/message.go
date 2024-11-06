package reusable

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
	userID := r.Context().Value(userIDKey).(int)

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

	query := `INSERT INTO Messages (user_id, content) VALUES ($1, $2) RETURNING id, created_at`
	err = db.QueryRow(query, userID, msg.Content).Scan(&msg.ID, &msg.CreatedAt)
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
			http.Error(w, "Не удалось получить сообщение", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	userID := r.Context().Value(userIDKey).(int)
	isPremium := r.Context().Value(isPremiumKey).(bool)

	if messageID == "" {
		http.Error(w, "Необходимо указать ID сообщения", http.StatusBadRequest)
		return
	}

	if !isPremium {
		http.Error(w, "Только премиум-пользователи могут редактировать свои сообщения", http.StatusForbidden)
		return
	}

	var updatedMsg Message
	err := json.NewDecoder(r.Body).Decode(&updatedMsg)
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

	var authorID int
	query := `SELECT user_id FROM Messages WHERE id = $1`
	err = db.QueryRow(query, messageID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Сообщение не найдено", http.StatusNotFound)
		} else {
			http.Error(w, "Не удалось получить данные сообщения", http.StatusInternalServerError)
		}
		return
	}

	if authorID != userID {
		http.Error(w, "Вы можете редактировать только свои сообщения", http.StatusForbidden)
		return
	}

	_, err = db.Exec(`UPDATE Messages SET content = $1 WHERE id = $2`, updatedMsg.Content, messageID)
	if err != nil {
		http.Error(w, "Не удалось обновить сообщение", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно обновлено"))
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	userID := r.Context().Value(userIDKey).(int)
	isPremium := r.Context().Value(isPremiumKey).(bool)

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

	var authorID int
	query := `SELECT user_id FROM Messages WHERE id = $1`
	err = db.QueryRow(query, messageID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Сообщение не найдено", http.StatusNotFound)
		} else {
			http.Error(w, "Не удалось получить сообщение", http.StatusInternalServerError)
		}
		return
	}

	if !isPremium || authorID != userID {
		http.Error(w, "Удаление разрешено только премиум-пользователям для своих сообщений", http.StatusForbidden)
		return
	}

	_, err = db.Exec(`DELETE FROM Messages WHERE id = $1`, messageID)
	if err != nil {
		http.Error(w, "Не удалось удалить сообщение", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно удалено"))
}

func LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	userID := r.Context().Value(userIDKey).(int)

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

	var likeExists bool
	query := `SELECT EXISTS(SELECT 1 FROM Likes WHERE user_id = $1 AND message_id = $2)`
	err = db.QueryRow(query, userID, messageID).Scan(&likeExists)
	if err != nil {
		http.Error(w, "Ошибка проверки лайка", http.StatusInternalServerError)
		return
	}

	if likeExists {
		http.Error(w, "Вы уже поставили лайк этому сообщению", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO Likes (user_id, message_id) VALUES ($1, $2)`, userID, messageID)
	if err != nil {
		http.Error(w, "Не удалось поставить лайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Лайк успешно добавлен"))
}

func SuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := r.URL.Query().Get("id")
	userID := r.Context().Value(userIDKey).(int)
	isPremium := r.Context().Value(isPremiumKey).(bool)

	if messageID == "" {
		http.Error(w, "Необходимо указать ID сообщения", http.StatusBadRequest)
		return
	}

	if !isPremium {
		http.Error(w, "Суперлайки доступны только премиум-пользователям", http.StatusForbidden)
		return
	}

	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var superlikeExists bool
	query := `SELECT EXISTS(SELECT 1 FROM SuperLikes WHERE user_id = $1 AND message_id = $2)`
	err = db.QueryRow(query, userID, messageID).Scan(&superlikeExists)
	if err != nil {
		http.Error(w, "Ошибка проверки суперлайка", http.StatusInternalServerError)
		return
	}

	if superlikeExists {
		http.Error(w, "Вы уже поставили суперлайк этому сообщению", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO SuperLikes (user_id, message_id) VALUES ($1, $2)`, userID, messageID)
	if err != nil {
		http.Error(w, "Не удалось поставить суперлайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Суперлайк успешно добавлен"))
}

func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, user_id, content, created_at FROM Messages`)
	if err != nil {
		http.Error(w, "Не удалось получить список сообщений", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err = rows.Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			http.Error(w, "Ошибка при обработке списка сообщений", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
