package redis

import (
    "context"
    "log"
    "os"
    "time"
    "github.com/redis/go-redis/v9"
    "github.com/Davi0805/gnose-notification/websocket"
    "github.com/Davi0805/gnose-notification/models"
)

type RedisService struct {
    client     *redis.Client
    hub        *websocket.Hub
    streamName string
}

func NewRedisService(hub *websocket.Hub) *RedisService {

    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }

    redisPassword := os.Getenv("REDIS_PASSWORD")

    
    client := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       0,
    })

    return &RedisService{
        client:     client,
        hub:        hub,
        streamName: "messages",
    }
}

func (s *RedisService) Setup(ctx context.Context) error {
    err := s.client.XGroupCreate(ctx, s.streamName, "message-group", "0").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        return err
    }
    return nil
}

func (s *RedisService) ConsumeMessages(ctx context.Context) {
    for {
        streams, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    "message-group",
            Consumer: "consumer-1",
            Streams:  []string{s.streamName, ">"},
        }).Result()
        if err != nil {
            log.Println("Erro na conexao com o redisStream:", err)
            time.Sleep(time.Second)
            continue
        }

        // TODO: REFAC
        for _, stream := range streams {
            for _, message := range stream.Messages {
                //id, ok := message.Values["id"]
                //if !ok {
                //    log.Println("Erro ao deserializar: id field nao encontrado")
                //    continue
                //}
//
                //idStr, ok := id.(string)
                //if !ok {
                //    log.Println("Erro ao deserializar: id field nao e uma string")
                //    continue
                //}



                content, ok := message.Values["content"]
                if !ok {
                    log.Println("Erro ao deserializar: content field nao foi encontrado")
                    continue
                }

                contentStr, ok := content.(string)
                if !ok {
                    log.Println("Erro ao deserializar: content field nao e uma string")
                    continue
                }

                companyId, ok := message.Values["companyId"]
                if !ok {
                    log.Println("Erro: companyId field nao foi encontrado")
                    continue
                }

                companyIdStr, ok := companyId.(string)
                if !ok {
                    log.Println("Erro: companyId field nao e uma string")
                    continue
                }

                userId, ok := message.Values["userId"]
                if !ok {
                    log.Println("Erro: companyId field nao foi encontrado")
                    continue
                }

                userIdStr, ok := userId.(string)
                if !ok {
                    log.Println("Erro: companyId field nao e uma string")
                    continue
                }

                service, ok := message.Values["service"]
                if !ok {
                    log.Println("Erro: service field nao foi encontrado")
                    continue
                }

                serviceStr, ok := service.(string)
                if !ok {
                    log.Println("Erro: service field nao e uma string")
                    continue
                }

                var msg models.Message
                msg.ID = message.ID
                msg.Content = contentStr
                msg.CompanyId = companyIdStr
                msg.UserId = userIdStr
                msg.Service = serviceStr
                s.hub.Broadcast(msg)
            }
        }
    }
}