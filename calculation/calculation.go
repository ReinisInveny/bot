package calculation

import (
	"github.com/markcheno/go-talib"
)

func ComputeStochRSI(prices []float64, rsi_period, k_period, d_period int) ([]float64, []float64, error) {

	fastk, fastd := talib.StochRsi(prices, rsi_period, k_period, d_period, talib.DEMA)

	return fastk, fastd, nil
}

func ComputeRSI(prices []float64, period int) ([]float64, error) {

	rsi := talib.Rsi(prices, period)

	return rsi, nil
}

func ComputeMACD(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) ([]float64, error) {

	_, _, macdHist := talib.Macd(closePrices, fastPeriod, slowPeriod, signalPeriod)

	return macdHist, nil
}
