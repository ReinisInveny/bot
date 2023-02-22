package database

import (
	"github.com/gangisreinis/bot/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() error {
	log.Println("bot connecting to database")

	conn, err := gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return err
	}

	DB = conn

	err = conn.AutoMigrate(&models.KlineData{}, &models.TechnicalIndicator{})

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
