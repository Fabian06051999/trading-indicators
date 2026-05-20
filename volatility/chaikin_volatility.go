package volatility

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
	"math"
)

// ChaikinVolatility implements Chaikin's Volatility indicator.
type ChaikinVolatility struct {
	emaPeriod int
	rocPeriod int
	ema       *moving_averages.EMA
	buffer    []float64
	index     int
	count     int
	out       []float64
}

func NewChaikinVolatility(emaPeriod, rocPeriod int) *ChaikinVolatility {
	return &ChaikinVolatility{
		emaPeriod: emaPeriod,
		rocPeriod: rocPeriod,
		ema:       moving_averages.NewEMA(emaPeriod),
		buffer:    make([]float64, rocPeriod+1),
		out:       make([]float64, 1),
	}
}

func (cv *ChaikinVolatility) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	cv.Reset()

	for i, c := range candles {
		cv.update(c)
		result[i] = cv.out[0]
	}
	return [][]float64{result}
}

func (cv *ChaikinVolatility) UpdateAll(candle indicators.OHLCV) []float64 {
	cv.update(candle)
	return cv.out
}

func (cv *ChaikinVolatility) update(candle indicators.OHLCV) {
	hl := candle.High - candle.Low
	emaVal := cv.ema.UpdateAll(indicators.OHLCV{Close: hl})[0]
	if math.IsNaN(emaVal) {
		cv.out[0] = math.NaN()
		return
	}

	cv.buffer[cv.index] = emaVal
	cv.count++
	cv.index = (cv.index + 1) % (cv.rocPeriod + 1)

	if cv.count <= cv.rocPeriod {
		cv.out[0] = math.NaN()
		return
	}

	pastVal := cv.buffer[cv.index]
	if pastVal == 0 {
		cv.out[0] = math.NaN()
		return
	}

	cv.out[0] = ((emaVal - pastVal) / pastVal) * 100
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
