package services

import (
	"github.com/gofiber/fiber/v2"
)

type GoldService interface {
	GetGold(c *fiber.Ctx) (*GoldService2, error)
	WebhookLineApi(c *fiber.Ctx) (*GoldService2, error)
}

type GoldService2 struct {
	Data     string `json:"data"`
	Messages string `json:"messages"`
}
