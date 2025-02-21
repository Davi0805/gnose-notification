package middleware

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "encoding/json"
    /* "github.com/golang-jwt/jwt/v4" */
    "github.com/redis/go-redis/v9"
	/* "log" */
    "fmt"
)

type TokenData struct {
    ID struct {
        UserId    int `json:"userId"`
        CompanyId int `json:"companyId"`
    } `json:"id"`
    Permission string `json:"permission"`
}

func JWTMiddleware(redisClient *redis.Client) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("Authorization")
        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Token malformatado"})
        }

		//PARSE DO TOKEN COM SIGN VERIFICATION
        /* token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Token invalido"})
        } */


        // SUBSTRING PARA REMOVER "Bearer" DO HEADER
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }

        ctx := context.Background()
        exists, err := redisClient.Exists(ctx, tokenString).Result()
        if err != nil || exists == 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Token invalido"})
        }

        // TODO: REFATORAR E LIMPAR PARA APENAS 1 QUERY QUANDO JA ESTIVER ESTAVEL
        // GET CREDENCIAIS DO REDIS
        jsonData, err := redisClient.Get(ctx, tokenString).Result()
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Erro ao buscar dados do token"})
        }

        // PARSE
        var tokenData []TokenData
        err = json.Unmarshal([]byte(jsonData), &tokenData)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Erro ao parsear dados do token"})
        }

        // CHECA SE EXISTE CREDENCIAL
        if len(tokenData) == 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Erro": "Dados do token nao encontrados"})
        }

        userId := tokenData[0].ID.UserId
        var companyIds []int
        for _, data := range tokenData {
            companyIds = append(companyIds, data.ID.CompanyId)
        }

        // DEBUG
        fmt.Println("Fetched userId:", userId)
        fmt.Println("Fetched companyId:", companyIds)

        // ENV VARIABLES
        c.Locals("userId", userId)
        c.Locals("companyIds", companyIds)

        return c.Next()
    }
}