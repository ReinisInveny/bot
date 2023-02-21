package database

import (
	"github.com/gangisreinis/bot/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	DB = conn

	err = conn.AutoMigrate(&models.KlineData{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot connected to database")

}
