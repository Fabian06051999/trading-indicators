package volatility

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// ChaikinVolatility implements Chaikin's Volatility indicator.
type ChaikinVolatility struct {
	emaPeriod int
	rocPeriod int
	ema       *moving_averages.EMA
	buffer    []float64
	index     int
	count     int
}

func NewChaikinVolatility(emaPeriod, rocPeriod int) *ChaikinVolatility {
	return &ChaikinVolatility{
		emaPeriod: emaPeriod,
		rocPeriod: rocPeriod,
		ema:       moving_averages.NewEMA(emaPeriod),
		buffer:    make([]float64, rocPeriod+1),
	}
}

func (cv *ChaikinVolatility) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	cv.Reset()

	for i, c := range candles {
		result[i] = cv.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (cv *ChaikinVolatility) UpdateAll(candle indicators.OHLCV) []float64 {
	hl := candle.High - candle.Low
	emaVal := cv.ema.UpdateAll(indicators.OHLCV{Close: hl})[0]
	if emaVal == 0 {
		return []float64{math.NaN()}
	}

	cv.buffer[cv.index] = emaVal
	cv.count++
	cv.index = (cv.index + 1) % (cv.rocPeriod + 1)

	if cv.count <= cv.rocPeriod {
		return []float64{math.NaN()}
	}

	pastVal := cv.buffer[cv.index]
	if pastVal == 0 {
		return []float64{math.NaN()}
	}

	return []float64{((emaVal - pastVal) / pastVal) * 100}
}

func (cv *ChaikinVolatility) Reset() {
	cv.ema.Reset()
	cv.buffer = make([]float64, cv.rocPeriod+1)
	cv.index = 0
	cv.count = 0
}

func (cv *ChaikinVolatility) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Chaikin Volatility",
		Parameters: []indicators.Parameter{
			{Name: "EMA Period", DefaultValue: float64(cv.emaPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "ROC Period", DefaultValue: float64(cv.rocPeriod), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "Chaikin Vol", Color: "#6A1B9A", Style: indicators.StyleLine, Width: 2},
		},
	}
}
