package reusable

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
	jwt.StandardClaims
}

type contextKey string

const (
	userIDKey    contextKey = "userID"
	isPremiumKey contextKey = "isPremium"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
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

	var storedPassword string
	var userID int
	var isPremium bool
	query := `SELECT id, password, is_premium FROM Users WHERE username = $1`
	err = db.QueryRow(query, creds.Username).Scan(&userID, &storedPassword, &isPremium)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		} else {
			http.Error(w, "Не удалось получить данные пользователя", http.StatusInternalServerError)
		}
		return
	}

	if creds.Password != storedPassword {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Не удалось создать токен", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		db, err := ConnectDB()
		if err != nil {
			http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var isPremium bool
		err = db.QueryRow(`SELECT is_premium FROM Users WHERE id = $1`, claims.UserID).Scan(&isPremium)
		if err != nil {
			http.Error(w, "Не удалось получить премиум-статус пользователя", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, isPremiumKey, isPremium)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
