package websocket

import (
    "github.com/Davi0805/gnose-notification/models"
    "github.com/Davi0805/gnose-notification/service"
    "github.com/gofiber/websocket/v2"
)

type Client struct {
    Conn *websocket.Conn
    User *models.User
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
            h.service.SaveMessage(message)
            for client := range h.clients {
                err := client.Conn.WriteJSON(message)
                if err != nil {
                    client.Conn.Close()
                    delete(h.clients, client)
                }
            }
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