package auth

import "github.com/dgrijalva/jwt-go"

type Role string

const (
	RoleUser    Role = "user"
	RolePremium Role = "premium"
)

type Claims struct {
	UserID int  `json:"user_id"`
	Role   Role `json:"role"`
	jwt.StandardClaims
}

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)
