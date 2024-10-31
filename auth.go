package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.StandardClaims
}

// Определение нового типа для ключей контекста
type contextKey string

const (
	userIDKey contextKey = "userID"
	roleIDKey contextKey = "roleID"
)

// Аутентификация пользователя и генерация JWT-токена
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
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

	// Проверка пользователя в базе данных
	var storedPassword string
	var userID, roleID int
	query := `SELECT id, password, role_id FROM Users WHERE username = $1`
	err = db.QueryRow(query, creds.Username).Scan(&userID, &storedPassword, &roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	// Сравнение паролей
	if creds.Password != storedPassword {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Создание JWT-токена
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		RoleID: roleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Возвращение токена клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Middleware для проверки токена JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавление данных о пользователе в контекст с использованием ключей userIDKey и roleIDKey
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, roleIDKey, claims.RoleID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Проверка роли суперпользователя
func IsSuperUser(r *http.Request) bool {
	roleID, ok := r.Context().Value(roleIDKey).(int)
	return ok && roleID == 2 // допустим, что роль "2" – это суперпользователь
}
