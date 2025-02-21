package repository

import (
    "database/sql"
    "github.com/Davi0805/gnose-notification/models"
    "fmt"
    "strings"
)

type MessageRepository struct {
    db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
    return &MessageRepository{db: db}
}

func (r *MessageRepository) Save(message models.Message) error {
    _, err := r.db.Exec("INSERT INTO messages (timestamp, content, company_id, user_id, service) VALUES ($1, $2, $3, $4, $5)",
        message.Timestamp, message.Content, message.CompanyId, message.UserId, message.Service)
    return err
}

func (r *MessageRepository) GetAll(companyIds []int64) ([]models.Message, error) {
    if len(companyIds) == 0 {
        return nil, nil
    }

    // Create placeholders for the query
    placeholders := make([]string, len(companyIds))
    args := make([]interface{}, len(companyIds))
    for i, id := range companyIds {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
        args[i] = id
    }

    query := fmt.Sprintf("SELECT id, timestamp, content, company_id, user_id, service FROM messages WHERE company_id IN (%s)", strings.Join(placeholders, ","))
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var message models.Message
        if err := rows.Scan(&message.ID, &message.Timestamp, &message.Content, &message.CompanyId, &message.UserId, &message.Service); err != nil {
            return nil, err
        }
        messages = append(messages, message)
    }
    return messages, nil
}

func (r *MessageRepository) GetByCompanyId(companyId int64) ([]models.Message, error) {
    rows, err := r.db.Query("SELECT id, timestamp, content, company_id, user_id, service FROM messages WHERE company_id = $1", companyId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []models.Message
    for rows.Next() {
        var message models.Message
        if err := rows.Scan(&message.ID, &message.Timestamp, &message.Content, &message.CompanyId, &message.UserId, &message.Service); err != nil {
            return nil, err
        }
        messages = append(messages, message)
    }
    return messages, nil
}