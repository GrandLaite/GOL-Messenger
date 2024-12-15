package users

import (
	"encoding/json"
	"gol_messenger/auth"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// Обработчик для регистрации нового пользователя
func (uh *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID, err := uh.UserService.RegisterUser(user)
	if err != nil {
		http.Error(w, "Не удалось зарегистрировать пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно зарегистрирован с ID: " + strconv.Itoa(userID)))
}

// Обработчик для логина пользователя
func (uh *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	token, err := uh.UserService.AuthenticateUser(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, "Ошибка авторизации: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Обработчик для получения информации о пользователе
func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok || userID == 0 {
		http.Error(w, "Не удалось получить ID пользователя", http.StatusUnauthorized)
		return
	}

	user, err := uh.UserService.GetUser(userID)
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Обработчик для обновления информации о пользователе
func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok || userID == 0 {
		http.Error(w, "Не удалось получить ID пользователя", http.StatusUnauthorized)
		return
	}

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	err = uh.UserService.UpdateUser(userID, updatedUser)
	if err != nil {
		http.Error(w, "Не удалось обновить информацию о пользователе", http.StatusInternalServerError)
		return
	}

	token, err := uh.UserService.GenerateToken(userID, updatedUser.Role)
	if err != nil {
		http.Error(w, "Не удалось создать новый токен", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Информация о пользователе успешно обновлена",
		"token":   token,
	})
}

// Обработчик для удаления пользователя
func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok || userID == 0 {
		http.Error(w, "Не удалось получить ID пользователя", http.StatusUnauthorized)
		return
	}

	err := uh.UserService.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Не удалось удалить пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пользователь успешно удалён"))
}
