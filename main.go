package main

import (
	"encoding/json"
	"github.com/gangisreinis/bot/binance"
	"github.com/gangisreinis/bot/calculation"
	"github.com/gangisreinis/bot/database"
	"github.com/gangisreinis/bot/models"
	"github.com/gangisreinis/bot/orm"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
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
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Panicln("could not connect to Binance websocket | Error:", err)
	}
	defer conn.Close()
	log.Println("websocket connection established")

	// Start a goroutine to read messages from the websocket
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Panicln("could not read message from websocket | Error:", err)
			}
			// Process the message here
			type KlineResponse struct {
				Value_e string           `json:"e"`
				Value_E int              `json:"E"`
				Value_s string           `json:"s"`
				Value_K models.KlineData `json:"k"`
			}

			var klineDataResp KlineResponse

			err = json.Unmarshal(message, &klineDataResp)
			if err != nil {
				log.Println("json error:", err)
				return
			}

			if klineDataResp.Value_K.KlineClosed {

				klines, err := orm.GetCandles(1000, interval, ticker)

				if err != nil {
					log.Println(err)
				}

				klines = append(klines, klineDataResp.Value_K)
				closePrices := make([]float64, 0, len(klines))
				lowPrices := make([]float64, 0, len(klines))
				highPrices := make([]float64, 0, len(klines))

				for _, kline := range klines {
					closePrice, err := strconv.ParseFloat(string(kline.ClosePrice), 64)
					if err != nil {
						log.Println("could not parse close price | kline:", kline, " | Error:", err)
						continue
					}
					lowPrice, err := strconv.ParseFloat(string(kline.LowPrice), 64)
					if err != nil {
						log.Println("could not parse low price | kline:", kline, " | Error:", err)
						continue
					}
					highPrice, err := strconv.ParseFloat(string(kline.HighPrice), 64)
					if err != nil {
						log.Println("could not parse high price | kline:", kline, " | Error:", err)
						continue
					}
					closePrices = append(closePrices, closePrice)
					lowPrices = append(lowPrices, lowPrice)
					highPrices = append(highPrices, highPrice)
				}

				RSI, err := calculation.ComputeRSI(closePrices, 14)
				if err != nil {
					log.Println(err)
				}

				MACD, err := calculation.ComputeMACD(closePrices, 12, 26, 9)
				if err != nil {
					log.Println(err)
				}

				// Clear the close prices slice for the next iteration
				closePrices = nil
				lowPrices = nil
				highPrices = nil

				// Wait for all goroutines to complete and collect their results

				var technicalIndicator models.TechnicalIndicator

				technicalIndicator.Ticker = ticker
				technicalIndicator.Interval = interval
				technicalIndicator.KlineCloseTime = klineDataResp.Value_K.KlineCloseTime
				technicalIndicator.ClosePrice, _ = klineDataResp.Value_K.ClosePrice.Float64()
				technicalIndicator.BaseAssetVolume, _ = klineDataResp.Value_K.BaseAssetVolume.Float64()
				technicalIndicator.NumberOfTrades = klineDataResp.Value_K.NumberOfTrades
				technicalIndicator.MACD = MACD[len(MACD)-1]
				technicalIndicator.RSI = RSI[len(RSI)-1]

				err = orm.AddCandle(klines[len(klines)-1])
				if err != nil {
					log.Println(err)
					return
				}

				err = orm.AddTechIndicator(technicalIndicator)
				if err != nil {
					log.Println(err)
					return
				}

				log.Printf("[WSS] %s | %.2f | MACD hist: %.2f | RSI: %.2f", ticker, technicalIndicator.ClosePrice, technicalIndicator.MACD, technicalIndicator.RSI)
			}
		}
	}()

	<-done
}
