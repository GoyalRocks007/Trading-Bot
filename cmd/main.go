package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"trading-bot/boot"
	"trading-bot/controllers"
	"trading-bot/db"
	"trading-bot/internal/bot"
	feedclient "trading-bot/internal/client/feed_client"
	"trading-bot/internal/models"
	executionmodule "trading-bot/internal/modules/execution_module"
	riskmodule "trading-bot/internal/modules/risk_module"
	tradestrategymodule "trading-bot/internal/modules/trade_strategy_module"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			panic("env file not found")
		}
	}
	router := gin.Default()
	db.InitDb()
	bus := models.NewBus(20000)
	// Register controller routes
	clientRegistery := boot.InitClientRegistery(feedclient.ZERODHA, bus, []string{os.Getenv("KAFKA_BROKER_URL")})
	adminController := controllers.NewAdminController(clientRegistery.FeedClient)

	adminRoutes := router.Group("/admin")
	{
		adminRoutes.GET("/auth", adminController.StartAuth)
	}
	// msg := map[string]interface{}{"to": "abhi25goyal@gmail.com", "Subject": "Testing Kafka.", "Body": "Kafka is working fine. Jai Shri Ram!"}
	// clientRegistery.Publisher.Publish(context.Background(), "notifications", "notification-group", msg)

	router.GET("/callback", adminController.HandleCallback)
	// router.GET("/hello", authController.HelloWorld)

	// // Start server
	// router.Run(":8080")
	wg := &sync.WaitGroup{}
	// ---------------------------------------------
	// START HTTP SERVER (BLOCKS)
	// ---------------------------------------------
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := router.Run(":8080")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Minute)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot := bot.NewBot(tradestrategymodule.ORB, executionmodule.PAPER, bus, riskmodule.SIMPLE, db.GetDbConnection(), 5*time.Minute)
		err := bot.Start()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	// Wait for both to finish
	wg.Wait()
}
