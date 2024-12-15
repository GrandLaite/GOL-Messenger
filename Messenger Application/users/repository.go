package users

import (
	"database/sql"
	"errors"
	"gol_messenger/auth"
)

type UserRepository interface {
	Create(user User) (int, error)
	GetByID(userID int) (User, error)
	Update(userID int, user User) error
	Delete(userID int) error
	GetByUsername(username string) (User, error)
	GetRoleByID(userID int) (auth.Role, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

// Создание пользователя
func (ur *userRepository) Create(user User) (int, error) {
	var userID int
	query := "INSERT INTO Users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id"
	err := ur.DB.QueryRow(query, user.Username, user.Password, user.Role).Scan(&userID)
	return userID, err
}

// Получение пользователя по ID
func (ur *userRepository) GetByID(userID int) (User, error) {
	var user User
	query := "SELECT id, username, password_hash, role FROM Users WHERE id = $1"
	err := ur.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return User{}, errors.New("пользователь не найден")
	}
	return user, err
}

// Обновление данных пользователя
func (ur *userRepository) Update(userID int, user User) error {
	query := "UPDATE Users SET username = $1, role = $2 WHERE id = $3"
	_, err := ur.DB.Exec(query, user.Username, user.Role, userID)
	return err
}

// Удаление пользователя
func (ur *userRepository) Delete(userID int) error {
	query := "DELETE FROM Users WHERE id = $1"
	_, err := ur.DB.Exec(query, userID)
	return err
}

// Получение пользователя по имени
func (ur *userRepository) GetByUsername(username string) (User, error) {
	var user User
	query := "SELECT id, username, password_hash, role FROM Users WHERE username = $1"
	err := ur.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return User{}, errors.New("пользователь не найден")
	}
	return user, err
}

// Получение роли пользователя по ID
func (ur *userRepository) GetRoleByID(userID int) (auth.Role, error) {
	var role auth.Role
	query := "SELECT role FROM Users WHERE id = $1"
	err := ur.DB.QueryRow(query, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return "", errors.New("роль пользователя не найдена")
	}
	return role, err
}
