package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/Davi0805/gnose-notification/service"
	"strconv"
)

type MessageController struct {
    service *service.MessageService
}

func NewMessageController(service *service.MessageService) *MessageController {
    return &MessageController{service: service}
}

func (c *MessageController) GetMessages(ctx *fiber.Ctx) error {
    messages, err := c.service.GetAllMessages()
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch messages",
        })
    }
    return ctx.JSON(messages)
}

func (c *MessageController) GetMessagesByCompanyId(ctx *fiber.Ctx) error {
    companyIdStr := ctx.Params("companyId")
    companyId, err := strconv.ParseInt(companyIdStr, 10, 64) // BASE 10 e 64 bits
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid company ID",
        })
    }
    messages, err := c.service.GetMessagesByCompanyId(companyId)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch messages",
        })
    }
    return ctx.JSON(messages)
}