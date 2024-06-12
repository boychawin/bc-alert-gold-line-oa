package configs

import (
	"bc-alert/src/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConfig() {

	viper.SetConfigType("json")
	viper.SetConfigName("environment." + os.Getenv("ENV"))
	viper.AddConfigPath("src/environments")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func InitTimeZone() {
	ict, err := time.LoadLocation(viper.GetString("TimeZone"))
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func InitCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*", // Single string origin
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type,Authorization,Token",
		AllowCredentials: false,
		// ExposeHeaders:    "Custom-Header",
	})
}

func FibersConfig() fiber.Config {
	return fiber.Config{
		BodyLimit: 200 * 1024 * 1024,
	}
}

// Use the LimitRequests middleware with a limit of 5 requests per minute
func LimitRequests(maxRequests int, duration time.Duration) fiber.Handler {
	requests := make(chan struct{}, maxRequests)
	ticker := time.Tick(duration)

	go func() {
		for range ticker {
			requests = make(chan struct{}, maxRequests)
		}
	}()

	return func(c *fiber.Ctx) error {
		select {
		case requests <- struct{}{}:
			return c.Next()
		default:
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"Status":   "error",
				"Messages": "Request limit exceeded",
			})
		}
	}
}

func InitDatabase() *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("DB.DB_USER"),
		viper.GetString("DB.DB_PASS"),
		viper.GetString("DB.DB_HOST"),
		viper.GetInt("DB.DB_POST"),
		viper.GetString("DB.DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &SqlLogger{}, 
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(models.BcAlertLine{}, models.BcAlertGolds{})

	fmt.Println("Database connection established successfully")

	return db
}
