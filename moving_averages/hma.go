package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// HMA implements the Hull Moving Average.
// HMA = WMA(2*WMA(n/2) - WMA(n), sqrt(n))
type HMA struct {
	period  int
	wmaHalf *WMA
	wmaFull *WMA
	wmaSqrt *WMA
	count   int
	out     []float64
}

func NewHMA(period int) *HMA {
	period = indicators.ClampMin(period, 2)
	sqrtPeriod := int(math.Sqrt(float64(period)))
	if sqrtPeriod < 1 {
		sqrtPeriod = 1
	}
	halfPeriod := period / 2
	if halfPeriod < 1 {
		halfPeriod = 1
	}

	return &HMA{
		period:  period,
		wmaHalf: NewWMA(halfPeriod),
		wmaFull: NewWMA(period),
		wmaSqrt: NewWMA(sqrtPeriod),
		out:     make([]float64, 1),
	}
}

func (h *HMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	h.Reset()

	for i, c := range candles {
		h.update(c)
		result[i] = h.out[0]
	}
	return [][]float64{result}
}

func (h *HMA) UpdateAll(candle indicators.OHLCV) []float64 {
	h.update(candle)
	return h.out
}

func (h *HMA) update(candle indicators.OHLCV) {
	h.count++
	halfVal := h.wmaHalf.UpdateAll(candle)[0]
	fullVal := h.wmaFull.UpdateAll(candle)[0]

	if math.IsNaN(halfVal) || math.IsNaN(fullVal) {
		h.out[0] = math.NaN()
		return
	}

	diff := 2*halfVal - fullVal
	h.out[0] = h.wmaSqrt.UpdateAll(indicators.OHLCV{Close: diff})[0]
}

func (h *HMA) Reset() {
	sqrtPeriod := int(math.Sqrt(float64(h.period)))
	if sqrtPeriod < 1 {
		sqrtPeriod = 1
	}
	halfPeriod := h.period / 2
	if halfPeriod < 1 {
		halfPeriod = 1
	}

	h.wmaHalf = NewWMA(halfPeriod)
	h.wmaFull = NewWMA(h.period)
	h.wmaSqrt = NewWMA(sqrtPeriod)
	h.count = 0
}

func (h *HMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Hull Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(h.period), Min: 2, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "HMA", Color: "#00BCD4", Style: indicators.StyleLine, Width: 2},
		},
	}
}
