package calculation

import (
	"github.com/markcheno/go-talib"
)

func ComputeStochRSI(high []float64, low []float64, close []float64, fastKP, slowKP, slowDP int) ([]float64, []float64, error) {

	// fastk, fastd := talib.StochRsi(high, low, close, fastKP, slowDP, talib.EMA, slowDP, talib.EMA)
	fastk, fastd := talib.StochRsi(close, 14, 3, 3, talib.DEMA)
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
