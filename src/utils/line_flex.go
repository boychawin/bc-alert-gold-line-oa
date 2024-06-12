package utils

import (
	"bc-alert/src/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func SendLineFlexMessage(to []string, request models.GoldPriceData) error {

	channelAccessToken := viper.GetString("LINE.CHANNEL_TOKEN")
	url := "https://api.line.me/v2/bot/message/multicast"

	flexMessage, err := BuildFlexContainer(request)
	if err != nil {
		log.Fatalf("Failed to generate Flex Container: %v", err)
	}

	toJSON, err := json.Marshal(to)
	if err != nil {
		log.Fatalf("Error marshaling 'to' slice:", err)
	}

	flexMessageJSON, err := json.Marshal(flexMessage)
	if err != nil {
		log.Fatalf("Failed to marshal Flex Message to JSON: %v", err)
	}

	payload := strings.NewReader(`{
		"to": ` + string(toJSON) + `,
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

	return err

}
