package middleware

import (
    "context"
    "github.com/gofiber/fiber/v2"
    /* "github.com/golang-jwt/jwt/v4" */
    "github.com/redis/go-redis/v9"
	/* "log" */
    /* "fmt" */
)

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

        return c.Next()
    }
}