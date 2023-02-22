package binance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gangisreinis/bot/calculation"
	"github.com/gangisreinis/bot/models"
)

func FetchKlines(ticker string, interval string) ([]models.KlineData, []models.TechnicalIndicator, error) {
	log.Println("requested binance API endpoint '/api/v3/klines/'")
	klineLimit := "1000"

	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%s", ticker, interval, klineLimit))
	if err != nil {
		return nil, nil, err
	}
	log.Println("received response from binance API")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var tempKlines [][]interface{}
	err = json.Unmarshal(body, &tempKlines)
	if err != nil {
		return nil, nil, err
	}

	var klines []models.KlineData
	var closePrices []float64
	// var lowPrices []float64
	// var highPrices []float64

	for i := 1; i < len(tempKlines)-1; i++ {
		if tempKlines[i][0].(float64) < tempKlines[i-1][0].(float64) {
			return nil, nil, fmt.Errorf("klines are not sorted in ascending order")
		}

		kline := models.KlineData{
			Ticker:                   ticker,
			Interval:                 interval,
			FirstTradeID:             0,
			LastTradeID:              0,
			KlineClosed:              true,
			KlineOpenTime:            int(tempKlines[i][0].(float64)),
			OpenPrice:                json.Number(tempKlines[i][1].(string)),
			HighPrice:                json.Number(tempKlines[i][2].(string)),
			LowPrice:                 json.Number(tempKlines[i][3].(string)),
			ClosePrice:               json.Number(tempKlines[i][4].(string)),
			BaseAssetVolume:          json.Number(tempKlines[i][5].(string)),
			KlineCloseTime:           int(tempKlines[i][6].(float64)),
			QuoteAssetVolume:         json.Number(tempKlines[i][7].(string)),
			NumberOfTrades:           int(tempKlines[i][8].(float64)),
			TakerBuyBaseAssetVolume:  json.Number(tempKlines[i][9].(string)),
			TakerBuyQuoteAssetVolume: json.Number(tempKlines[i][10].(string)),
		}
		klines = append(klines, kline)

		closePrice, err := strconv.ParseFloat(string(kline.ClosePrice), 64)
		if err != nil {
			return nil, nil, err
		}
		// lowPrice, err := strconv.ParseFloat(string(kline.LowPrice), 64)
		// if err != nil {
		// 	log.Println("could not parse low price | kline:", kline, " | Error:", err)
		// 	continue
		// }
		// highPrice, err := strconv.ParseFloat(string(kline.HighPrice), 64)
		// if err != nil {
		// 	log.Println("could not parse high price | kline:", kline, " | Error:", err)
		// 	continue
		// }
		closePrices = append(closePrices, closePrice)
		// lowPrices = append(lowPrices, lowPrice)
		// highPrices = append(highPrices, highPrice)
	}

	// Create channels to receive the results of the goroutines
	ch1 := make(chan []float64)
	ch2 := make(chan []float64)
	defer close(ch1)
	defer close(ch2)

	// Define the goroutines

	go func() {
		RSI, err := calculation.ComputeRSI(closePrices, 14)
		if err != nil {
			log.Println(err)
		}
		ch1 <- RSI
	}()

	go func() {
		MACD, err := calculation.ComputeMACD(closePrices, 12, 26, 9)
		if err != nil {
			log.Println(err)
		}
		ch2 <- MACD
	}()

	// Wait for all goroutines to complete and collect their results
	RSI := <-ch1
	MACD := <-ch2

	// Combine the results into the final output
	var technicalIndicators []models.TechnicalIndicator
	var technicalIndicator models.TechnicalIndicator

	for i, kline := range klines {
		technicalIndicator.Ticker = ticker
		technicalIndicator.Interval = interval
		technicalIndicator.KlineCloseTime = kline.KlineCloseTime
		technicalIndicator.ClosePrice, _ = kline.ClosePrice.Float64()
		technicalIndicator.BaseAssetVolume, _ = kline.BaseAssetVolume.Float64()
		technicalIndicator.NumberOfTrades = kline.NumberOfTrades
		technicalIndicator.MACD = MACD[i]
		technicalIndicator.RSI = RSI[i]

		technicalIndicators = append(technicalIndicators, technicalIndicator)
	}

	return klines, technicalIndicators, nil
}
