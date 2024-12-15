package messages

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gol_messenger/auth"
)

type MessageHandler struct {
	MessageService MessageService
}

func NewMessageHandler(messageService MessageService) *MessageHandler {
	return &MessageHandler{MessageService: messageService}
}

func (mh *MessageHandler) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	role, ok := r.Context().Value(auth.UserRoleKey).(auth.Role)
	if !ok || role != auth.RoleUser {
		http.Error(w, "Только пользователи могут создавать сообщения", http.StatusForbidden)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	createdMsg, err := mh.MessageService.CreateMessage(userID, msg.Content)
	if err != nil {
		http.Error(w, "Не удалось создать сообщение", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdMsg)
}

func (mh *MessageHandler) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageIDStr := r.URL.Query().Get("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Неверный ID сообщения", http.StatusBadRequest)
		return
	}

	msg, err := mh.MessageService.GetMessage(messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (mh *MessageHandler) LikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	role, ok := r.Context().Value(auth.UserRoleKey).(auth.Role)
	if !ok || role != auth.RoleUser {
		http.Error(w, "Только пользователи могут ставить лайки", http.StatusForbidden)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Неверный ID сообщения", http.StatusBadRequest)
		return
	}

	err = mh.MessageService.LikeMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Не удалось поставить лайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Лайк поставлен!"))
}

func (mh *MessageHandler) SuperlikeMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	role, ok := r.Context().Value(auth.UserRoleKey).(auth.Role)
	if !ok || role != auth.RolePremium {
		http.Error(w, "Только премиум-пользователи могут ставить супер-лайки", http.StatusForbidden)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Неверный ID сообщения", http.StatusBadRequest)
		return
	}

	err = mh.MessageService.SuperlikeMessage(messageID, userID)
	if err != nil {
		http.Error(w, "Не удалось поставить супер-лайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Супер-лайк поставлен!"))
}

func (mh *MessageHandler) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	role, ok := r.Context().Value(auth.UserRoleKey).(auth.Role)
	if !ok || role != auth.RolePremium {
		http.Error(w, "Удаление сообщений доступно только премиум-пользователям", http.StatusForbidden)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Неверный ID сообщения", http.StatusBadRequest)
		return
	}

	err = mh.MessageService.DeleteMessage(messageID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно удалено"))
}

func (mh *MessageHandler) UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(int)
	role, ok := r.Context().Value(auth.UserRoleKey).(auth.Role)
	if !ok || role != auth.RolePremium {
		http.Error(w, "Обновление сообщений доступно только пользователям", http.StatusForbidden)
		return
	}

	messageIDStr := r.URL.Query().Get("id")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		http.Error(w, "Неверный ID сообщения", http.StatusBadRequest)
		return
	}

	var msg Message
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	err = mh.MessageService.UpdateMessage(messageID, userID, msg.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно обновлено"))
}

func (mh *MessageHandler) ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := mh.MessageService.ListMessages()
	if err != nil {
		http.Error(w, "Не удалось получить сообщения", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
