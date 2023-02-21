package main

import (
	"github.com/gangisreinis/bot/binance"
	"github.com/gangisreinis/bot/database"
	"github.com/gangisreinis/bot/orm"
	"log"
)

func main() {
	log.Println("staring up the bot")
	ticker := "BTCEUR"
	interval := "1m"
	log.Println("selected ticker:", ticker, " | selected interval:", interval)
	database.Connect()
	log.Println("initilized connection to database")

	// Clear database of historic data

	if err := orm.DeleteCandles(); err != nil {
		log.Panicln("could not clear historic candles Error:", err)
	}

	// Load new candle data

	klines, err := binance.FetchKlines(ticker, interval)

	if err != nil {
		log.Panicln("could not fetch candles | Error", err)
	}

	for _, kline := range klines {
		if err := orm.AddCandle(kline); err != nil {
			log.Panicln("could not add kline to table | kline:", kline, " | Error:", err)
		}
	}

	log.Println("new klines are added to database")

	log.Println("starting up websocket connection")

}
