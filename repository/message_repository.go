package repository

import "github.com/Davi0805/gnose-notification/models"

type MessageRepository struct {
    messages   []models.Message
    saveChan   chan models.Message
    getAllChan chan chan []models.Message
}

func NewMessageRepository() *MessageRepository {
    repo := &MessageRepository{
        messages:   make([]models.Message, 0),
        saveChan:   make(chan models.Message),
        getAllChan: make(chan chan []models.Message),
    }
    go repo.run()
    return repo
}

func (r *MessageRepository) run() {
    for {
        select {
        case message := <-r.saveChan:
            r.messages = append(r.messages, message)
        case replyChan := <-r.getAllChan:
            replyChan <- r.messages
        }
    }
}

func (r *MessageRepository) Save(message models.Message) error {
    r.saveChan <- message
    return nil
}

func (r *MessageRepository) GetAll() ([]models.Message, error) {
    replyChan := make(chan []models.Message)
    r.getAllChan <- replyChan
    messages := <-replyChan
    return messages, nil
}