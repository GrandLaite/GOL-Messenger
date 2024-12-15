package users

import "gol_messenger/auth"

// User - структура пользователя
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     auth.Role `json:"role"`
}
