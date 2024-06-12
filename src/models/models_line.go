package models

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type Event struct {
	Type            string `json:"type"`
	Timestamp       int64  `json:"timestamp"`
	Source          Source `json:"source"`
	ReplyToken      string `json:"replyToken"`
	Mode            string `json:"mode"`
	WebhookEventID  string `json:"webhookEventId"`
	DeliveryContext struct {
		IsRedelivery bool `json:"isRedelivery"`
	} `json:"deliveryContext"`
}

type LineWebhookPayload struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}
