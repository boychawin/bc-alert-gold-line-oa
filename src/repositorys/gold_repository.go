package repositorys

import (
	"bc-alert/src/models"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type goldRepositoryDB struct {
	db *gorm.DB
}

func NewGoldRepositoryDB(db *gorm.DB) GoldRepository {
	return goldRepositoryDB{db}
}

func (u goldRepositoryDB) GetGold(c *fiber.Ctx) (*models.GoldPriceResponse, error) {

	var (
		barSell       string
		barBuy        string
		ornamentSell  string
		ornamentBuy   string
		updateTheTime string
		todayChange   string
		updatedDate   string
		updatedTime   string
	)

	c_colly := colly.NewCollector()

	c_colly.OnHTML(viper.GetString("GOLD.GoldtradersTableElement"), func(e *colly.HTMLElement) {
		barSell = e.ChildText("tr:nth-child(1) > td:nth-child(2)")
		barBuy = e.ChildText("tr:nth-child(1) > td:nth-child(3)")
		ornamentSell = e.ChildText("tr:nth-child(2) > td:nth-child(2)")
		ornamentBuy = e.ChildText("tr:nth-child(2) > td:nth-child(3)")
		// statusChange = e.ChildAttr("tr:nth-child(3) > td.span > .imgp", "alt")
		todayChange = strings.Split(e.ChildText("tr:nth-child(3) > td:nth-child(1)"), " ")[1]
		updatedDate = e.ChildText("tr:nth-child(4) > td.span.bg-span.txtd.al-r")
		updatedTime = e.ChildText("tr:nth-child(4) > td.em.bg-span.txtd.al-r")
		updateTheTime = e.ChildText("tr:nth-child(4) > td.em.bg-span.txtd.al-l")
	})

	c_colly.Visit(viper.GetString("GOLD.GoldtradersLink"))

	response := &models.GoldPriceResponse{
		ResponseData: models.GoldPriceData{
			BarBuy:        barBuy,
			BarSell:       barSell,
			OrnamentBuy:   ornamentBuy,
			OrnamentSell:  ornamentSell,
			StatusChange:  "-50",
			TodayChange:   todayChange,
			UpdatedDate:   updatedDate,
			UpdatedTime:   updatedTime,
			UpdateTheTime: updateTheTime,
		},
	}

	return response, nil

}

func (u goldRepositoryDB) WebhookLineApi(c *fiber.Ctx) (*models.GoldPriceResponse, error) {
	var request = &models.LineWebhookPayload{}

	origin := c.Get("Origin")
	log.Println("Origin :", origin)

	if err := c.BodyParser(&request); err != nil {
		return nil, errors.New("request is nil after parsing body")
	}

	if request == nil {
		return nil, errors.New("request is nil after parsing body")
	}

	if len(request.Events) == 0 {
		return nil, errors.New("no events in the request")
	}

	event := request.Events[0]

	bcAlertLine := models.BcAlertLine{
		Type:           event.Type,
		Timestamp:      event.Timestamp,
		UserID:         event.Source.UserID,
		ReplyToken:     event.ReplyToken,
		Mode:           event.Mode,
		WebhookEventID: event.WebhookEventID,
		IsRedelivery:   event.DeliveryContext.IsRedelivery,
	}

	existingRecord := models.BcAlertLine{}
	if err := u.db.First(&existingRecord, "user_id = ?", bcAlertLine.UserID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if existingRecord.ID == 0 {

		resultCreate := u.db.Create(&bcAlertLine)
		if resultCreate.Error != nil {
			fmt.Println("Error:", resultCreate.Error)
			return nil, resultCreate.Error
		}
	} else {

		existingRecord.Type = bcAlertLine.Type
		existingRecord.Timestamp = bcAlertLine.Timestamp
		existingRecord.ReplyToken = bcAlertLine.ReplyToken
		existingRecord.Mode = bcAlertLine.Mode
		existingRecord.WebhookEventID = bcAlertLine.WebhookEventID
		existingRecord.IsRedelivery = bcAlertLine.IsRedelivery

		resultUpdate := u.db.Save(&existingRecord)
		if resultUpdate.Error != nil {
			fmt.Println("Error:", resultUpdate.Error)
			return nil, resultUpdate.Error
		}
	}

	response := &models.GoldPriceResponse{
		ResponseMessage: "success",
	}

	return response, nil
}
