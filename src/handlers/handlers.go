package handlers

import (

	"github.com/gofiber/fiber/v2"

)

type GoldHandler interface {
	GetGold(c *fiber.Ctx) error
	WebhookLineApi(c *fiber.Ctx) error
}
