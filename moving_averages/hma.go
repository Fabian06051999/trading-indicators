package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// HMA implements the Hull Moving Average.
// HMA = WMA(2*WMA(n/2) - WMA(n), sqrt(n))
type HMA struct {
	period    int
	wmaHalf   *WMA
	wmaFull   *WMA
	wmaSqrt   *WMA
	count     int
}

func NewHMA(period int) *HMA {
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
	}
}

func (h *HMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	h.Reset()

	for i, c := range candles {
		result[i] = h.Update(c)
	}
	return result
}

func (h *HMA) Update(candle indicators.OHLCV) float64 {
	h.count++
	halfVal := h.wmaHalf.Update(candle)
	fullVal := h.wmaFull.Update(candle)

	if halfVal == 0 || fullVal == 0 {
		return 0
	}

	diff := 2*halfVal - fullVal
	result := h.wmaSqrt.Update(indicators.OHLCV{Close: diff})
	return result
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
