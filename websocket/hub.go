package websocket

import (
    "github.com/Davi0805/gnose-notification/models"
    "github.com/Davi0805/gnose-notification/service"
    "github.com/gofiber/websocket/v2"
    "strconv"
)

type Client struct {
    Conn *websocket.Conn
    User *models.User
}

// TODO: REFATORAR E ALTERAR DB SCHEME PARA ARMAZENAR IDS COMO LONG E EVITAR ATOI
func (c *Client) IsPartOfCompany(companyId string) bool {
    companyIdInt, err := strconv.Atoi(companyId)
    if err != nil {
        return false
    }
    for _, id := range c.User.CompanyIds {
        if id == companyIdInt {
            return true
        }
    }
    return false
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan models.Message
    register   chan *Client
    unregister chan *Client
    service    *service.MessageService
}

func NewHub(service *service.MessageService) *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan models.Message),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        service:    service,
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                client.Conn.Close()
            }
        case message := <-h.broadcast:
            go h.runBroadcast(message)
        }
    }
}

func (h *Hub) Register(client *Client) {
    h.register <- client
}

func (h *Hub) Unregister(client *Client) {
    h.unregister <- client
}

// BROADCAST PARA TODOS OS CLIENTS
func (h *Hub) Broadcast(message models.Message) {
    h.broadcast <- message
}

func (h *Hub) runBroadcast(message models.Message) {
    h.service.SaveMessage(message)
    for client := range h.clients {
        if client.IsPartOfCompany(message.CompanyId) {
            err := client.Conn.WriteJSON(message)
            if err != nil {
                client.Conn.Close()
                delete(h.clients, client)
            }
        }
    }
}