package main

import (
	"bc-alert/src/configs"
	"bc-alert/src/handlers"
	"bc-alert/src/models"
	"bc-alert/src/repositorys"
	"bc-alert/src/services"
	"bc-alert/src/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/gocolly/colly"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func main() {
	configs.InitTimeZone()
	configs.InitConfig()
	db := configs.InitDatabase()

	app := fiber.New(configs.FibersConfig())

	app.Use(configs.InitCors())
	app.Use(configs.LimitRequests(1000, time.Minute))

	goldRepositoryDB := repositorys.NewGoldRepositoryDB(db)
	goldService := services.NewGoldService(goldRepositoryDB)
	goldHandler := handlers.NewGoldHandler(goldService)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, BC-Alert!")
	})

	app.Post("/api/line/webhook", goldHandler.WebhookLineApi)

	// Set up cron job
	cronSchedule := "*/30 * * * * *"
	// cronSchedule := "* * * * **
	job := cron.New()
	job.AddFunc(cronSchedule, func() {
		GetGoldCron(db)

	})
	go job.Start()

	// Start the Fiber server
	err := app.Listen(fmt.Sprintf(":%v", viper.GetInt("LISTENING.PORT")))
	if err != nil {
		panic(err)
	}
}

func GetGoldCron(db *gorm.DB) error {

	isCheckDB := true
	ENV_MODE := viper.GetString("ENV")
	if ENV_MODE == "dev" {
		isCheckDB = true
	}
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

	if isCheckDB {
		barSellFloat, _ := strconv.ParseFloat(strings.Replace(barSell, ",", "", -1), 64)
		barBuyFloat, _ := strconv.ParseFloat(strings.Replace(barBuy, ",", "", -1), 64)
		ornamentSellFloat, _ := strconv.ParseFloat(strings.Replace(ornamentSell, ",", "", -1), 64)
		ornamentBuyFloat, _ := strconv.ParseFloat(strings.Replace(ornamentBuy, ",", "", -1), 64)
		todayChangeFloat, _ := strconv.ParseFloat(strings.Replace(todayChange, ",", "", -1), 64)

		bcAlert := models.BcAlertGolds{
			BarBuy:        barBuyFloat,
			BarSell:       barSellFloat,
			OrnamentBuy:   ornamentBuyFloat,
			OrnamentSell:  ornamentSellFloat,
			TodayChange:   todayChangeFloat,
			StatusChange:  "",
			UpdatedDate:   updatedDate,
			UpdatedTime:   updatedTime,
			UpdateTheTime: updateTheTime,
		}

		if bcAlert.UpdatedTime == "" || bcAlert.UpdatedDate == "" {
			fmt.Println("Alert:UpdatedTime is NULL")
			return nil
		}

		var existingRecord models.BcAlertGolds
		result := db.Where("updated_time LIKE ? AND updated_date LIKE ?", bcAlert.UpdatedTime, bcAlert.UpdatedDate).First(&existingRecord)
		if result.Error == nil {
			fmt.Println("Alert:exists")
			return result.Error
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {

			fmt.Println("Error:", result.Error)
			return result.Error
		}

		resultCreate := db.Create(&bcAlert)
		if resultCreate.Error != nil {
			fmt.Println("Error:", resultCreate.Error)
		} else {
			fmt.Println("Alert: saved successfully")
		}

		data_string := "ทองคำ\n"
		data_string += fmt.Sprintf("ซื้อบาท: %.2f\n", barBuyFloat)
		data_string += fmt.Sprintf("ขายบาท: %.2f\n", barSellFloat)
		data_string += fmt.Sprintf("ซื้อทองรูปพรรณ: %.2f\n", ornamentBuyFloat)
		data_string += fmt.Sprintf("ขายทองรูปพรรณ: %.2f\n", ornamentSellFloat)
		data_string += fmt.Sprintf("เปลี่ยนแปลง: %.2f\n", todayChangeFloat)
		data_string += fmt.Sprintf("อัพเดท: %s/%s\n", updatedDate, updatedTime)

		utils.SendMessageToLineNotify(data_string)

	}

	var bcAlertLines []models.BcAlertLine
	result := db.Where("user_id IS NOT NULL AND status_paid = TRUE").Find(&bcAlertLines)
	if result.Error != nil {
		fmt.Println("Alert: line exists")
		return result.Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		fmt.Println("Error: line", result.Error)
		return result.Error
	}

	if len(bcAlertLines) < 0 {
		fmt.Println("Alert: bcAlertLine < 0")
		return nil
	}
	var to []string
	for _, line := range bcAlertLines {
		to = append(to, line.UserID)
	}

	err := utils.SendLineFlexMessage(to, models.GoldPriceData{
		BarBuy:        barBuy,
		BarSell:       barSell,
		OrnamentBuy:   ornamentBuy,
		OrnamentSell:  ornamentSell,
		StatusChange:  "-50",
		TodayChange:   todayChange,
		UpdatedDate:   updatedDate,
		UpdatedTime:   updatedTime,
		UpdateTheTime: updateTheTime,
	})
	if err != nil {
		return err
	}

	return nil

}
