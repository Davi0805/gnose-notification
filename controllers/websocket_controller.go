package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/Davi0805/gnose-notification/models"
    ws "github.com/Davi0805/gnose-notification/websocket"
)

type WebSocketController struct {
    hub *ws.Hub
}

func NewWebSocketController(hub *ws.Hub) *WebSocketController {
    return &WebSocketController{hub: hub}
}

func (c *WebSocketController) HandleWebSocket(ctx *fiber.Ctx) error {
    return websocket.New(func(conn *websocket.Conn) {
        client := &ws.Client{Conn: conn}
        c.hub.Register(client)

        defer func() {
            c.hub.Unregister(client)
        }()

        for {
            var message models.Message
            err := conn.ReadJSON(&message)
            if err != nil {
                break
            }
            c.hub.Broadcast(message)
        }
    })(ctx)
}