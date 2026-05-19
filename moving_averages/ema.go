package moving_averages

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// EMA implements the Exponential Moving Average.
type EMA struct {
	period     int
	multiplier float64
	value      float64
	count      int
	sum        float64
}

func NewEMA(period int) *EMA {
	return &EMA{
		period:     period,
		multiplier: 2.0 / float64(period+1),
	}
}

func (e *EMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	e.Reset()

	for i, c := range candles {
		result[i] = e.Update(c)
	}
	return result
}

func (e *EMA) Update(candle indicators.OHLCV) float64 {
	e.count++
	if e.count <= e.period {
		e.sum += candle.Close
		if e.count == e.period {
			e.value = e.sum / float64(e.period)
		}
		return math.NaN()
	}
	e.value = (candle.Close-e.value)*e.multiplier + e.value
	return e.value
}

func (e *EMA) Reset() {
	e.value = 0
	e.count = 0
	e.sum = 0
}

func (e *EMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Exponential Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(e.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "EMA", Color: "#FF9800", Style: indicators.StyleLine, Width: 2},
		},
	}
}
