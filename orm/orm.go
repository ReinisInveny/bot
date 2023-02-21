package orm

import (
	"fmt"
	"github.com/gangisreinis/bot/database"
	"github.com/gangisreinis/bot/models"
	"time"
)

func AddCandle(candle models.KlineData) error {

	if err := database.DB.Create(&candle).Error; err != nil {
		return err
	}

	return nil
}

func DeleteCandles() error {

	if err := database.DB.Exec("DELETE FROM kline_data").Error; err != nil {
		return err
	}

	return nil
}

func GetCandles(limit int, interval, ticker string) ([]models.KlineData, error) {
	endTime := time.Now().UnixNano() / 1000000

	startTime := int64(0)

	switch interval {
	case "1m":
		startTime = endTime - (int64(limit) * 60000)

	case "3m":
		startTime = endTime - (int64(limit) * 60000 * 3)

	case "5m":
		startTime = endTime - (int64(limit) * 60000 * 5)

	case "15m":
		startTime = endTime - (int64(limit) * 60000 * 15)

	case "30m":
		startTime = endTime - (int64(limit) * 60000 * 30)

	case "1h":
		startTime = endTime - (int64(limit) * 60000 * 60)

	case "2h":
		startTime = endTime - (int64(limit) * 60000 * 120)

	case "4h":
		startTime = endTime - (int64(limit) * 60000 * 240)

	case "6h":
		startTime = endTime - (int64(limit) * 60000 * 360)

	case "8h":
		startTime = endTime - (int64(limit) * 60000 * 480)

	case "12h":
		startTime = endTime - (int64(limit) * 60000 * 720)

	case "1d":
		startTime = endTime - (int64(limit) * 60000 * 1440)

	case "3d":
		startTime = endTime - (int64(limit) * 60000 * 4320)

	case "1w":
		startTime = endTime - (int64(limit) * 60000 * 10080)

	case "1M":
		startTime = endTime - (int64(limit) * 60000 * 43800)
	}

	var candles []models.KlineData
	err := database.DB.Where("ticker = ? AND interval = ? AND kline_close_time >= ? AND kline_close_time <= ?", ticker, interval, startTime, endTime).Find(&candles).Error
	if err != nil {
		return nil, err
	}

	if len(candles) != limit {
		return nil, fmt.Errorf("number of candles returned does not match the specified limit value")
	}

	return candles, nil
}
