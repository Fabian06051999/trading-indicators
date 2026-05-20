package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// WilliamsAD implements Williams Accumulation/Distribution.
type WilliamsAD struct {
	out       []float64
	value     float64
	prevClose float64
	count     int
}

func NewWilliamsAD() *WilliamsAD {
	return &WilliamsAD{out: make([]float64, 1)}
}

func (w *WilliamsAD) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		result[i] = w.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (w *WilliamsAD) UpdateAll(candle indicators.OHLCV) []float64 {
	w.count++

	if w.count == 1 {
		w.prevClose = candle.Close
		w.out[0] = math.NaN()
		return w.out
	}

	ad := 0.0
	if candle.Close > w.prevClose {
		ad = candle.Close - min(candle.Low, w.prevClose)
	} else if candle.Close < w.prevClose {
		ad = candle.Close - max(candle.High, w.prevClose)
	}

	w.value += ad
	w.prevClose = candle.Close
	w.out[0] = w.value
	return w.out
}

func (w *WilliamsAD) Reset() {
	w.value = 0
	w.prevClose = 0
	w.count = 0
}

func (w *WilliamsAD) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name:       "Williams Accumulation/Distribution",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "WAD", Color: "#4527A0", Style: indicators.StyleLine, Width: 2},
		},
	}
}
