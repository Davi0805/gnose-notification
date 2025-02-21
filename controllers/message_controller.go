package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/Davi0805/gnose-notification/service"
	"strconv"
	"context"
)

type MessageController struct {
    service *service.MessageService
	authService *service.AuthService
}

func NewMessageController(service *service.MessageService, authService *service.AuthService) *MessageController {
    return &MessageController{service: service, authService: authService}
}

func (c *MessageController) GetMessages(ctx *fiber.Ctx) error {

	// GET HEADER
    authHeader := ctx.Get("Authorization")
    if authHeader == "" {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "ERRO": "Authorization header nao encontrado",
        })
    }

    // GET CREDENTIALS
    credentials, err := c.authService.GetCredentialsFromToken(context.Background(), authHeader)
    if err != nil {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "ERRO": err.Error(),
        })
    }

	// GET EMPRESAS DO USUARIO
    companyIds := c.getCompanyIds(credentials)


    messages, err := c.service.GetAllMessages(companyIds)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "ERRO": "Mensagem nao encontrada",
        })
    }
    return ctx.JSON(messages)
}

func (c *MessageController) GetMessagesByCompanyId(ctx *fiber.Ctx) error {


    companyIdStr := ctx.Params("companyId")
    companyId, err := strconv.ParseInt(companyIdStr, 10, 64) // BASE 10 e 64 bits
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "ERRO": "Company id invalido",
        })
    }

	// GET HEADER
    authHeader := ctx.Get("Authorization")
    if authHeader == "" {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "ERRO": "Authorization header nao encontrado",
        })
    }

    // GET CREDENTIALS
    credentials, err := c.authService.GetCredentialsFromToken(context.Background(), authHeader)
    if err != nil {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "ERRO": err.Error(),
        })
    }

    // CHECA SE FAZ PARTE DA EMPRESA
    isAuthorized := false
    for _, cred := range credentials {
        if cred.ID.CompanyId == int(companyId) {
            isAuthorized = true
            break
        }
    }

    if !isAuthorized {
        return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "ERRO": "User nao faz parte da empresa",
        })
    }

    messages, err := c.service.GetMessagesByCompanyId(companyId)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "ERRO": "Mensagem nao encontrada",
        })
    }
    return ctx.JSON(messages)
}

func (c *MessageController) getCompanyIds(credentials []service.TokenData) []int64 {
    var companyIds []int64
    for _, cred := range credentials {
        companyIds = append(companyIds, int64(cred.ID.CompanyId))
    }
    return companyIds
}