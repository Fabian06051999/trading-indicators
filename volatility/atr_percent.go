package volatility

import (
	"math"
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

func (a *ATRPercent) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		result[i] = a.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (a *ATRPercent) UpdateAll(candle indicators.OHLCV) []float64 {
	atrVal := a.atr.UpdateAll(candle)[0]
	if atrVal == 0 || candle.Close == 0 {
		return []float64{math.NaN()}
	}
	return []float64{(atrVal / candle.Close) * 100}
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
