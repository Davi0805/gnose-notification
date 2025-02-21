package service

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/redis/go-redis/v9"
)

type TokenData struct {
    ID struct {
        UserId    int `json:"userId"`
        CompanyId int `json:"companyId"`
    } `json:"id"`
    Permission string `json:"permission"`
}

type AuthService struct {
    redisClient *redis.Client
}

func NewAuthService(redisClient *redis.Client) *AuthService {
    return &AuthService{redisClient: redisClient}
}

func (s *AuthService) GetCredentialsFromToken(ctx context.Context, tokenString string) ([]TokenData, error) {
    
	// REMOVE "BEARER"
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }


    // GET CREDENCIAIS
    jsonData, err := s.redisClient.Get(ctx, tokenString).Result()
    if err != nil {
        return nil, errors.New("Falha no sistema de credentials")
    }

    // PARSE
    var tokenData []TokenData
    err = json.Unmarshal([]byte(jsonData), &tokenData)
    if err != nil {
        return nil, errors.New("Falha de parse nas credenciais")
    }

    return tokenData, nil
}