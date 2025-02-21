package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/Davi0805/gnose-notification/models"
    ws "github.com/Davi0805/gnose-notification/websocket"
    /* "fmt" */
)

type WebSocketController struct {
    hub *ws.Hub
}

func NewWebSocketController(hub *ws.Hub) *WebSocketController {
    return &WebSocketController{hub: hub}
}

func (c *WebSocketController) HandleWebSocket(ctx *fiber.Ctx) error {
    // GET CONTEXT VARIABLES ANTES DE DAR UPGRADE DE CONTEXTO
    userId, ok := ctx.Locals("userId").(int)
    if !ok || userId == 0 {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User ID not found",
        })
    }

    companyIds, ok := ctx.Locals("companyIds").([]int)
    if !ok || len(companyIds) == 0 {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Company IDs not found",
        })
    }

    // SETANDO USUARIO COM VARIAVEIS DA AUTH
    return websocket.New(func(conn *websocket.Conn) {
        client := &ws.Client{
            Conn: conn,
            User: &models.User{
                ID:         userId,
                CompanyIds: companyIds,
            },
        }

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
