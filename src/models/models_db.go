package models

import (
	"gorm.io/gorm"
)

type BcAlertGolds struct {
	gorm.Model
	BarBuy        float64 `gorm:"type:decimal(10,2)"`
	BarSell       float64 `gorm:"type:decimal(10,2)"`
	OrnamentBuy   float64 `gorm:"type:decimal(10,2)"`
	OrnamentSell  float64 `gorm:"type:decimal(10,2)"`
	StatusChange  string  `gorm:"type:varchar(255)"`
	TodayChange   float64 `gorm:"type:decimal(10,2)"`
	UpdatedDate   string  `gorm:"type:varchar(255)"`
	UpdatedTime   string  `gorm:"type:varchar(255)"`
	UpdateTheTime string  `gorm:"type:varchar(255)"`
}

type BcAlertLine struct {
	gorm.Model
	DisplayName    string `gorm:"type:text"`
	PictureURL     string `gorm:"type:text"`
	AccessToken    string `gorm:"type:text"`
	Language       string `gorm:"type:text"`
	OS             string `gorm:"type:text"`
	LineVersion    string `gorm:"type:text"`
	LoggedIn       bool   `gorm:"default:false"`
	UserID         string `gorm:"type:text"`
	ReplyToken     string `gorm:"type:text"`
	WebhookEventID string `gorm:"type:text"`
	Mode           string `gorm:"type:text"`
	Timestamp      int64
	Type           string
	IsRedelivery   bool
	StatusPaid     bool `gorm:"default:false"`
}
