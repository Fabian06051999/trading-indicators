package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// ATR implements the Average True Range.
type ATR struct {
	period    int
	value     float64
	prevClose float64
	count     int
	sum       float64
}

func NewATR(period int) *ATR {
	return &ATR{
		period: period,
	}
}

func (a *ATR) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		result[i] = a.Update(c)
	}
	return result
}

func (a *ATR) Update(candle indicators.OHLCV) float64 {
	a.count++

	if a.count == 1 {
		a.prevClose = candle.Close
		tr := candle.High - candle.Low
		a.sum = tr
		return math.NaN()
	}

	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-a.prevClose), math.Abs(candle.Low-a.prevClose)))
	a.prevClose = candle.Close

	if a.count <= a.period {
		a.sum += tr
		if a.count == a.period {
			a.value = a.sum / float64(a.period)
			return a.value
		}
		return math.NaN()
	}

	// Wilder's smoothing
	a.value = (a.value*float64(a.period-1) + tr) / float64(a.period)
	return a.value
}

func (a *ATR) Reset() {
	a.value = 0
	a.prevClose = 0
	a.count = 0
	a.sum = 0
}

func (a *ATR) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Average True Range",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(a.period), Min: 1, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "ATR", Color: "#FF5722", Style: indicators.StyleLine, Width: 2},
		},
	}
}
