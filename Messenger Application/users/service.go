package users

import (
	"errors"
	"gol_messenger/auth"
)

type UserService interface {
	AuthenticateUser(username, password string) (string, error)
	GetUser(userID int) (User, error)
	UpdateUser(userID int, updatedUser User) error
	DeleteUser(userID int) error
	RegisterUser(user User) (int, error)
	GenerateToken(userID int, role auth.Role) (string, error)
}

type userService struct {
	UserRepository UserRepository
	TokenService   auth.TokenService
}

func NewUserService(userRepository UserRepository, tokenService auth.TokenService) UserService {
	return &userService{
		UserRepository: userRepository,
		TokenService:   tokenService,
	}
}

// Регистрация пользователя
func (us *userService) RegisterUser(user User) (int, error) {
	// Хэширование пароля
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return 0, errors.New("ошибка хэширования пароля")
	}
	user.Password = hashedPassword

	// Установка роли по умолчанию
	if user.Role == "" {
		user.Role = auth.RoleUser
	}

	// Создание пользователя в репозитории
	userID, err := us.UserRepository.Create(user)
	if err != nil {
		return 0, errors.New("ошибка создания пользователя")
	}

	return userID, nil
}

// Аутентификация пользователя
func (us *userService) AuthenticateUser(username, password string) (string, error) {
	// Получение пользователя по имени
	user, err := us.UserRepository.GetByUsername(username)
	if err != nil {
		return "", errors.New("неверное имя пользователя или пароль")
	}

	// Проверка пароля
	if !auth.CheckPasswordHash(password, user.Password) {
		return "", errors.New("неверное имя пользователя или пароль")
	}

	// Генерация токена
	token, err := us.TokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("ошибка генерации токена")
	}

	return token, nil
}

// Генерация токена для пользователя
func (us *userService) GenerateToken(userID int, role auth.Role) (string, error) {
	token, err := us.TokenService.GenerateToken(userID, role)
	if err != nil {
		return "", errors.New("ошибка генерации токена")
	}
	return token, nil
}

// Получение информации о пользователе
func (us *userService) GetUser(userID int) (User, error) {
	return us.UserRepository.GetByID(userID)
}

// Обновление данных пользователя
func (us *userService) UpdateUser(userID int, updatedUser User) error {
	// Хэширование пароля при обновлении
	if updatedUser.Password != "" {
		hashedPassword, err := auth.HashPassword(updatedUser.Password)
		if err != nil {
			return errors.New("ошибка хэширования пароля")
		}
		updatedUser.Password = hashedPassword
	}

	return us.UserRepository.Update(userID, updatedUser)
}

// Удаление пользователя
func (us *userService) DeleteUser(userID int) error {
	return us.UserRepository.Delete(userID)
}
