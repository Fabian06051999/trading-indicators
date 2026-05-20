package volatility

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// ATRPercent implements the ATR as a percentage of price (normalized ATR).
type ATRPercent struct {
	atr *ATR
	out []float64
}

func NewATRPercent(period int) *ATRPercent {
	return &ATRPercent{
		atr: NewATR(period),
		out: make([]float64, 1),
	}
}

func (a *ATRPercent) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		a.update(c)
		result[i] = a.out[0]
	}
	return [][]float64{result}
}

func (a *ATRPercent) UpdateAll(candle indicators.OHLCV) []float64 {
	a.update(candle)
	return a.out
}

func (a *ATRPercent) update(candle indicators.OHLCV) {
	atrVal := a.atr.UpdateAll(candle)[0]
	if math.IsNaN(atrVal) || candle.Close == 0 {
		a.out[0] = math.NaN()
		return
	}
	a.out[0] = (atrVal / candle.Close) * 100
}

func (a *ATRPercent) Reset() {
	a.atr.Reset()
}

func (a *ATRPercent) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "ATR Percent",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(a.atr.period), Min: 1, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "ATR%", Color: "#E65100", Style: indicators.StyleLine, Width: 2},
		},
	}
}
