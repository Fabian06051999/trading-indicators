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
	out       []float64
}

func NewATR(period int) *ATR {
	period = indicators.ClampMin(period, 1)
	return &ATR{
		period: period,
		out:    make([]float64, 1),
	}
}

func (a *ATR) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		a.update(c)
		result[i] = a.out[0]
	}
	return [][]float64{result}
}

func (a *ATR) UpdateAll(candle indicators.OHLCV) []float64 {
	a.update(candle)
	return a.out
}

func (a *ATR) update(candle indicators.OHLCV) {
	a.count++

	if a.count == 1 {
		a.prevClose = candle.Close
		tr := candle.High - candle.Low
		a.sum = tr
		a.out[0] = math.NaN()
		return
	}

	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-a.prevClose), math.Abs(candle.Low-a.prevClose)))
	a.prevClose = candle.Close

	if a.count <= a.period {
		a.sum += tr
		if a.count == a.period {
			a.value = a.sum / float64(a.period)
			a.out[0] = a.value
			return
		}
		a.out[0] = math.NaN()
		return
	}

	// Wilder's smoothing
	a.value = (a.value*float64(a.period-1) + tr) / float64(a.period)
	a.out[0] = a.value
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
