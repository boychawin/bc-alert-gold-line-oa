package services

import (
	"bc-alert/src/repositorys"
	"bc-alert/src/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type goldService struct {
	goldRepo repositorys.GoldRepository
}

func NewGoldService(goldRepo repositorys.GoldRepository) GoldService {
	return goldService{goldRepo}
}

func (s goldService) GetGold(c *fiber.Ctx) (*GoldService2, error) {
	request, err := s.goldRepo.GetGold(c)
	if err != nil {
		return nil, nil
	}

	channelAccessToken := viper.GetString("LINE.CHANNEL_TOKEN")
	url := "https://api.line.me/v2/bot/message/push"
	to := "U65579bc26af0d67c4ef90dbdaeaf98ec"

	flexMessage, err := utils.BuildFlexContainer(request.ResponseData)
	if err != nil {
		log.Fatalf("Failed to generate Flex Container: %v", err)
	}

	flexMessageJSON, err := json.Marshal(flexMessage)
	if err != nil {
		log.Fatalf("Failed to marshal Flex Message to JSON: %v", err)
	}

	payload := strings.NewReader(`{
		"to": "` + to + `",
		"messages": [` + string(flexMessageJSON) + `]
	}`)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Failed to send message, status code: %d", res.StatusCode)
	}

	log.Println("Message sent successfully!")

	response := &GoldService2{
		Data:     "",
		Messages: "success",
	}

	return response, nil
}

func (s goldService) WebhookLineApi(c *fiber.Ctx) (*GoldService2, error) {
	request, err := s.goldRepo.WebhookLineApi(c)
	if err != nil {
		return nil, nil
	}

	response := &GoldService2{
		Data:     "",
		Messages: request.ResponseMessage,
	}

	return response, nil
}
