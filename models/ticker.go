package models

import (
	"encoding/json"
)

type KlineData struct {
	ID                       uint        `gorm:"primaryKey"`
	KlineOpenTime            int         `json:"t"`
	KlineCloseTime           int         `json:"T"`
	Ticker                   string      `json:"s"`
	Interval                 string      `json:"i"`
	FirstTradeID             int         `json:"f"`
	LastTradeID              int         `json:"L"`
	OpenPrice                json.Number `json:"o"`
	ClosePrice               json.Number `json:"c"`
	HighPrice                json.Number `json:"h"`
	LowPrice                 json.Number `json:"l"`
	BaseAssetVolume          json.Number `json:"v"`
	NumberOfTrades           int         `json:"n"`
	KlineClosed              bool        `json:"x"`
	QuoteAssetVolume         json.Number `json:"q"`
	TakerBuyBaseAssetVolume  json.Number `json:"V"`
	TakerBuyQuoteAssetVolume json.Number `json:"Q"`
}

type TechnicalIndicator struct {
	ID              uint `gorm:"primaryKey"`
	KlineCloseTime  int
	Ticker          string
	Interval        string
	ClosePrice      float64
	BaseAssetVolume float64
	NumberOfTrades  int
	MACD            float64
	RSI             float64
}
