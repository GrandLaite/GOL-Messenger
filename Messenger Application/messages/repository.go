package messages

import (
	"database/sql"
	"errors"
)

type MessageRepository interface {
	Create(userID int, content string) (Message, error)
	GetByID(id int) (Message, error)
	GetAll() ([]Message, error)
	Update(id int, content string) error
	Delete(id int) error
	InsertLike(id int, userID int) error
	InsertSuperlike(id int, userID int) error
}

type messageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{DB: db}
}

func (mr *messageRepository) Create(userID int, content string) (Message, error) {
	var createdMessage Message
	query := `INSERT INTO Messages (user_id, content) VALUES ($1, $2) RETURNING id, user_id, content, created_at`
	err := mr.DB.QueryRow(query, userID, content).Scan(
		&createdMessage.ID,
		&createdMessage.UserID,
		&createdMessage.Content,
		&createdMessage.CreatedAt,
	)
	return createdMessage, err
}

func (mr *messageRepository) GetByID(id int) (Message, error) {
	var message Message
	query := `SELECT id, user_id, content, created_at FROM Messages WHERE id = $1`
	err := mr.DB.QueryRow(query, id).Scan(
		&message.ID,
		&message.UserID,
		&message.Content,
		&message.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return Message{}, errors.New("сообщение не найдено")
	}
	return message, err
}

func (mr *messageRepository) GetAll() ([]Message, error) {
	query := `SELECT id, user_id, content, created_at FROM Messages`
	rows, err := mr.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.UserID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (mr *messageRepository) Update(id int, content string) error {
	query := `UPDATE Messages SET content = $1 WHERE id = $2`
	_, err := mr.DB.Exec(query, content, id)
	return err
}

func (mr *messageRepository) Delete(id int) error {
	query := `DELETE FROM Messages WHERE id = $1`
	_, err := mr.DB.Exec(query, id)
	return err
}

func (mr *messageRepository) InsertLike(id int, userID int) error {
	query := `INSERT INTO Likes (message_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := mr.DB.Exec(query, id, userID)
	return err
}

func (mr *messageRepository) InsertSuperlike(id int, userID int) error {
	query := `INSERT INTO SuperLikes (message_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := mr.DB.Exec(query, id, userID)
	return err
}
