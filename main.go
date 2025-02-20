package main

import (
    "context"
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/Davi0805/gnose-notification/repository"
    "github.com/Davi0805/gnose-notification/service"
    ws "github.com/Davi0805/gnose-notification/websocket"
    "github.com/Davi0805/gnose-notification/redis"
    "github.com/Davi0805/gnose-notification/controllers"
)

func main() {
    app := fiber.New()

    // INIT DAS DEPENDENCIAS
    repo := repository.NewMessageRepository()
    service := service.NewMessageService(repo)
    hub := ws.NewHub(service)
    controller := controllers.NewWebSocketController(hub)

    // INICIA HUB DO WEB SOCKET
    go hub.Run()

    // INICIALIZA REDIS NO BACKGROUND
    redisService := redis.NewRedisService(hub)
    ctx := context.Background()
    if err := redisService.Setup(ctx); err != nil {
        log.Fatalf("Failed to setup Redis service: %v", err)
    }
    go redisService.ConsumeMessages(ctx)

    // TODO: ADICIONAR AUTH
    // DECLARA ENDPOINT DO WS
    app.Use("/ws", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    // TODO: ROTAS HTTP COM NOTIFICACOES LIGADO A DB
    app.Get("/ws", controller.HandleWebSocket)

    log.Fatal(app.Listen(":3000"))
}