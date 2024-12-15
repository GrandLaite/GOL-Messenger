package messages

import (
	"errors"
)

type MessageService interface {
	CreateMessage(userID int, content string) (Message, error)
	GetMessage(messageID int) (Message, error)
	ListMessages() ([]Message, error)
	UpdateMessage(messageID int, userID int, content string) error
	DeleteMessage(messageID int, userID int) error
	LikeMessage(messageID int, userID int) error
	SuperlikeMessage(messageID int, userID int) error
}

type messageService struct {
	Repository MessageRepository
}

func NewMessageService(repository MessageRepository) MessageService {
	return &messageService{Repository: repository}
}

func (ms *messageService) CreateMessage(userID int, content string) (Message, error) {
	return ms.Repository.Create(userID, content)
}

func (ms *messageService) GetMessage(messageID int) (Message, error) {
	return ms.Repository.GetByID(messageID)
}

func (ms *messageService) ListMessages() ([]Message, error) {
	return ms.Repository.GetAll()
}

func (ms *messageService) UpdateMessage(messageID int, userID int, content string) error {
	msg, err := ms.Repository.GetByID(messageID)
	if err != nil {
		return err
	}

	if msg.UserID != userID {
		return errors.New("операция запрещена")
	}

	return ms.Repository.Update(messageID, content)
}

func (ms *messageService) DeleteMessage(messageID int, userID int) error {
	msg, err := ms.Repository.GetByID(messageID)
	if err != nil {
		return err
	}

	if msg.UserID != userID {
		return errors.New("операция запрещена")
	}

	return ms.Repository.Delete(messageID)
}

func (ms *messageService) LikeMessage(messageID int, userID int) error {
	return ms.Repository.InsertLike(messageID, userID)
}

func (ms *messageService) SuperlikeMessage(messageID int, userID int) error {
	return ms.Repository.InsertSuperlike(messageID, userID)
}
