package main

import (
	"log"
	"strings"
	"time"

	"github.com/gangisreinis/bot/binance"
	"github.com/gangisreinis/bot/database"
	"github.com/gangisreinis/bot/orm"
)

const (
	ticker   = "BTCEUR"
	interval = "1m"
)

func main() {
	log.Println("starting up the bot")
	addr := "wss://stream.binance.com:9443/ws/" + strings.ToLower(ticker) + "@kline_" + interval
	log.Println("selected ticker:", ticker, " | selected interval:", interval)

	// Connect to the database
	if err := database.Connect(); err != nil {
		log.Panicln("could not connect to database | Error:", err)
	}
	log.Println("database connection established")

	// Clear historic data from the database
	if err := orm.DeleteCandles(); err != nil {
		log.Panicln("could not clear historic candles | Error:", err)
	}
	if err := orm.DeleteTechInd(); err != nil {
		log.Panicln("could not clear historic technical indicators | Error:", err)
	}

	// Fetch new candle data and add it to the database
	klines, techIndicators, err := binance.FetchKlines(ticker, interval)
	if err != nil {
		log.Panicln("could not fetch candles | Error", err)
	}

	if err := orm.AddCandles(klines); err != nil {
		log.Panicln("could not add candles to database | Error:", err)
	}
	log.Printf("%d new klines added to database\n", len(klines))

	if err := orm.AddTechIndicators(techIndicators); err != nil {
		log.Panicln("could not add technical indicators to database | Error:", err)
	}
	log.Printf("%d new technical indicators added to database\n", len(techIndicators))

	// Start up the websocket connection

	update := make(chan bool)
	done := make(chan bool)

	go binance.Websocket(update, done, interval, ticker, addr)
	for {
		select {
		case <-update:
			onUpdate()
		case <-time.After(30 * time.Second):
			onWait()
		}
	}

}

func onUpdate() {
	log.Println("update state")
}

func onWait() {
	log.Println("waiting state")
}
