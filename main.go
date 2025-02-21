package main

import (
    "context"
    "log"
    "os"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/Davi0805/gnose-notification/repository"
    "github.com/Davi0805/gnose-notification/service"
    ws "github.com/Davi0805/gnose-notification/websocket"
    "github.com/Davi0805/gnose-notification/redis"
    "github.com/Davi0805/gnose-notification/controllers"
    "github.com/Davi0805/gnose-notification/middleware"
    "github.com/joho/godotenv"
    redisv9 "github.com/redis/go-redis/v9"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    app := fiber.New()

    // POSTGRES CREDENTIALS
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // ! SE N ME ENGANO EM C sprintf n e memory safe mas aq n deve ter problema
    dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    // INIT DAS DEPENDENCIAS
    db, err := repository.NewPostgresDB(dataSource)
    if err != nil {
        log.Fatalf("Falha ao conectar ao db: %v", err)
    }
    defer db.Close()

    // TODO: TALVEZ SUBSTITUIR POR MESMO CLIENTE USADO NO WS | TESTAR EM TESTE DE CARGA
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    redisPassword := os.Getenv("REDIS_PASSWORD")
    redisClient := redisv9.NewClient(&redisv9.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       0,
    })

    authService := service.NewAuthService(redisClient)

    repo := repository.NewMessageRepository(db)
    messageService := service.NewMessageService(repo)
    hub := ws.NewHub(messageService)
    controller := controllers.NewWebSocketController(hub)
    messageController := controllers.NewMessageController(messageService, authService)

    // INICIA HUB DO WEB SOCKET
    go hub.Run()

    // INICIALIZA REDIS NO BACKGROUND
    redisService := redis.NewRedisService(hub)
    ctx := context.Background()
    if err := redisService.Setup(ctx); err != nil {
        log.Fatalf("Failed to setup Redis service: %v", err)
    }
    go redisService.ConsumeMessages(ctx)

    // DECLARA ENDPOINT DO WS
    app.Use("/ws", middleware.JWTMiddleware(redisService.GetClient()), func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    // TODO: ROTAS HTTP COM NOTIFICACOES LIGADO A DB
    app.Get("/ws", controller.HandleWebSocket)
    
    app.Get("/messages", messageController.GetMessages)
    app.Get("/messages/:companyId", messageController.GetMessagesByCompanyId)

    log.Fatal(app.Listen(":3000"))
}