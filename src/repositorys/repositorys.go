package repositorys

import (
	"bc-alert/src/models"

	"github.com/gofiber/fiber/v2"
)

type GoldRepository interface {
	GetGold(c *fiber.Ctx) (*models.GoldPriceResponse, error)
	WebhookLineApi(c *fiber.Ctx) (*models.GoldPriceResponse, error)
}
