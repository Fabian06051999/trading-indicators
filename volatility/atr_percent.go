package volatility

import (
	"github.com/Fabian06051999/trading-indicators"
)

// ATRPercent implements the ATR as a percentage of price (normalized ATR).
type ATRPercent struct {
	atr *ATR
}

func NewATRPercent(period int) *ATRPercent {
	return &ATRPercent{
		atr: NewATR(period),
	}
}

func (a *ATRPercent) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		result[i] = a.Update(c)
	}
	return result
}

func (a *ATRPercent) Update(candle indicators.OHLCV) float64 {
	atrVal := a.atr.Update(candle)
	if atrVal == 0 || candle.Close == 0 {
		return 0
	}
	return (atrVal / candle.Close) * 100
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
