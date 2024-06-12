package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func SendMessageToLineNotify(message string) error {
	accessToken := viper.GetString("LINE_GROUP")

	formData := url.Values{}
	formData.Set("message", message)

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to LINE Notify. Status: %s", resp.Status)
	}

	return nil
}
