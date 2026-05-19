package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// WilliamsAD implements Williams Accumulation/Distribution.
type WilliamsAD struct {
	value     float64
	prevClose float64
	count     int
}

func NewWilliamsAD() *WilliamsAD {
	return &WilliamsAD{}
}

func (w *WilliamsAD) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		result[i] = w.Update(c)
	}
	return result
}

func (w *WilliamsAD) Update(candle indicators.OHLCV) float64 {
	w.count++

	if w.count == 1 {
		w.prevClose = candle.Close
		return math.NaN()
	}

	ad := 0.0
	if candle.Close > w.prevClose {
		ad = candle.Close - min(candle.Low, w.prevClose)
	} else if candle.Close < w.prevClose {
		ad = candle.Close - max(candle.High, w.prevClose)
	}

	w.value += ad
	w.prevClose = candle.Close
	return w.value
}

func (w *WilliamsAD) Reset() {
	w.value = 0
	w.prevClose = 0
	w.count = 0
}

func (w *WilliamsAD) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Williams Accumulation/Distribution",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "WAD", Color: "#4527A0", Style: indicators.StyleLine, Width: 2},
		},
	}
}

