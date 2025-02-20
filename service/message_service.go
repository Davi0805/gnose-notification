package service

import (
    "github.com/Davi0805/gnose-notification/models"
    "github.com/Davi0805/gnose-notification/repository"
)

type MessageService struct {
    repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
    return &MessageService{repo: repo}
}

func (s *MessageService) SaveMessage(message models.Message) error {
    return s.repo.Save(message)
}

func (s *MessageService) GetAllMessages() ([]models.Message, error) {
    return s.repo.GetAll()
}