package handlers

import (
	"bc-alert/src/models"
	"bc-alert/src/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type goldHandler struct {
	goldSrv services.GoldService
}

func NewGoldHandler(goldSrv services.GoldService) GoldHandler {
	return goldHandler{goldSrv}
}


func (h goldHandler) GetGold(c *fiber.Ctx) error {

	req, err := h.goldSrv.GetGold(c)

	if err != nil {
		return c.JSON(fiber.Map{
			"Status":  "error",
			"Messages": err.Error(),
		})
	}

	response := fiber.Map{
		"Status":   "ok",
		"Data":       req.Data,
		"Messages":   req.Messages,
	}

	return c.JSON(response)

}

func (h goldHandler) WebhookLineApi(c *fiber.Ctx) error {

	if err := c.Body(); len(err) == 0 {
		fmt.Println(err)
        return  c.Status(fiber.StatusOK).JSON(fiber.Map{
            "Status":   "error",
            "Messages": err,
        })
    }

	var request = &models.LineWebhookPayload{}

	if err := c.BodyParser(&request); err != nil {
		return  c.Status(fiber.StatusOK).JSON(fiber.Map{
            "Status":   "error",
            "Messages": "events is empty",
        })
	}

	

	if request.Events == nil {
        return  c.Status(fiber.StatusOK).JSON(fiber.Map{
            "Status":   "error",
            "Messages": "events is empty",
        })
    }

    // ตรวจสอบจำนวน items ใน events
    if len(request.Events) == 0 {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "Status":   "error",
            "Messages": "events must contain at least 1 item",
        })
    }



	req, err := h.goldSrv.WebhookLineApi(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Status":  "error",
			"Messages": err.Error(),
		})
	}

	response := fiber.Map{
		"Status":   "ok",
		"Data":       req.Data,
		"Messages":   req.Messages,
	}

	return c.JSON(response)

}

