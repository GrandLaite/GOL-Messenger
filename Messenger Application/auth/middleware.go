package auth

import (
	"context"
	"net/http"
)

type AuthMiddleware struct {
	TokenService TokenService
}

func NewAuthMiddleware(tokenService TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		TokenService: tokenService,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
			return
		}

		claims, err := am.TokenService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) RoleMiddleware(requiredRole Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
				return
			}

			claims, err := am.TokenService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Неверный токен", http.StatusUnauthorized)
				return
			}

			if claims.Role != requiredRole {
				http.Error(w, "Доступ запрещён: недостаточно прав", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
