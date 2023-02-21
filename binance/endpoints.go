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

func FetchKlines(ticker string, interval string) ([]models.KlineData, error) {
	log.Println("requested binance API endpoint '/api/v3/klines/'")
	klineLimit := "1000"

	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%s", ticker, interval, klineLimit))
	if err != nil {
		return nil, err
	}
	log.Println("received response from binance API")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tempKlines [][]interface{}
	err = json.Unmarshal(body, &tempKlines)
	if err != nil {
		return nil, err
	}

	var klines []models.KlineData
	var closePrices []float64

	for i := 1; i < len(tempKlines); i++ {
		if tempKlines[i][0].(float64) < tempKlines[i-1][0].(float64) {
			return nil, fmt.Errorf("klines are not sorted in ascending order")
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
			return nil, err
		}
		closePrices = append(closePrices, closePrice)
	}

	// Create channels to receive the results of the goroutines
	ch1 := make(chan []float64)
	ch2 := make(chan []float64)
	ch3 := make(chan []float64)
	ch4 := make(chan []float64)
	defer close(ch1)
	defer close(ch2)
	defer close(ch3)
	defer close(ch4)

	// Define the goroutines
	go func() {
		fastK, fastD, err := calculation.ComputeStochRSI(closePrices, 14, 3, 3)
		if err != nil {
			log.Println(err)
		}
		ch1 <- fastK
		ch2 <- fastD
	}()

	go func() {
		RSI, err := calculation.ComputeRSI(closePrices, 14)
		if err != nil {
			log.Println(err)
		}
		ch3 <- RSI
	}()

	go func() {
		MACD, err := calculation.ComputeMACD(closePrices, 12, 26, 9)
		if err != nil {
			log.Println(err)
		}
		ch4 <- MACD
	}()

	// Wait for all goroutines to complete and collect their results
	fastK := <-ch1
	fastD := <-ch2
	RSI := <-ch3
	MACD := <-ch4

	// Combine the results into the final output
	for i := range klines {
		klines[i].MACD = MACD[i]
		klines[i].RSI = RSI[i]
		klines[i].RSI_STOCH_FAST_K = fastK[i]
		klines[i].RSI_STOCH_FAST_D = fastD[i]
	}

	return klines, nil
}
