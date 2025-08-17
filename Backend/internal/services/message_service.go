package services

import (
	"github.com/rafael-bit/whatz/internal/models"
	"github.com/rafael-bit/whatz/internal/repository"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
}

func NewMessageService(messageRepo *repository.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) Create(message *models.Message) error {
	return s.messageRepo.Create(message)
}

func (s *MessageService) GetByID(id string) (*models.Message, error) {
	return s.messageRepo.GetByID(id)
}

func (s *MessageService) GetByRoom(roomID string, limit, offset int) ([]*models.Message, error) {
	return s.messageRepo.GetByRoom(roomID, limit, offset)
}

func (s *MessageService) GetRecentMessages(roomID string, limit int) ([]*models.Message, error) {
	return s.messageRepo.GetRecentMessages(roomID, limit)
}

func (s *MessageService) GetByUser(userID string, limit, offset int) ([]*models.Message, error) {
	return s.messageRepo.GetByUser(userID, limit, offset)
}

func (s *MessageService) Delete(id string) error {
	return s.messageRepo.Delete(id)
}

func (s *MessageService) GetMessageCount(roomID string) (int, error) {
	return s.messageRepo.GetMessageCount(roomID)
}
