package db

import (
	"log"
	"os"
	"trading-bot/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DbConnection    *gorm.DB
	GetDbConnection = func() *gorm.DB {
		if DbConnection != nil {
			return DbConnection
		}
		return nil
	}
)

func InitDb() {
	var dsn string = os.Getenv("DB_CONNECTION_URL")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DbConnection = db

	err = db.AutoMigrate(&models.Order{}, &models.Position{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("Database migrated successfully 🚀")
}
