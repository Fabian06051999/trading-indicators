package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// WMA implements the Weighted Moving Average.
type WMA struct {
	period int
	buffer []float64
	index  int
	count  int
	out    []float64
}

func NewWMA(period int) *WMA {
	period = indicators.ClampMin(period, 1)
	return &WMA{
		period: period,
		buffer: make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (w *WMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		w.update(c)
		result[i] = w.out[0]
	}
	return [][]float64{result}
}

func (w *WMA) UpdateAll(candle indicators.OHLCV) []float64 {
	w.update(candle)
	return w.out
}

func (w *WMA) update(candle indicators.OHLCV) {
	w.buffer[w.index] = candle.Close
	w.index = (w.index + 1) % w.period
	if w.count < w.period {
		w.count++
	}
	if w.count < w.period {
		w.out[0] = math.NaN()
		return
	}

	weightSum := 0.0
	divisor := 0.0
	for i := 0; i < w.period; i++ {
		idx := (w.index + i) % w.period
		weight := float64(i + 1)
		weightSum += w.buffer[idx] * weight
		divisor += weight
	}
	w.out[0] = weightSum / divisor
}

func (w *WMA) Reset() {
	w.buffer = make([]float64, w.period)
	w.index = 0
	w.count = 0
}

func (w *WMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Weighted Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(w.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "WMA", Color: "#4CAF50", Style: indicators.StyleLine, Width: 2},
		},
	}
}
