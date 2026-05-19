package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
)

// WilliamsR implements Williams %R.
type WilliamsR struct {
	period int
	highs  []float64
	lows   []float64
	index  int
	count  int
}

func NewWilliamsR(period int) *WilliamsR {
	return &WilliamsR{
		period: period,
		highs:  make([]float64, period),
		lows:   make([]float64, period),
	}
}

func (w *WilliamsR) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		result[i] = w.Update(c)
	}
	return result
}

func (w *WilliamsR) Update(candle indicators.OHLCV) float64 {
	w.highs[w.index] = candle.High
	w.lows[w.index] = candle.Low
	w.index = (w.index + 1) % w.period
	if w.count < w.period {
		w.count++
	}
	if w.count < w.period {
		return 0
	}

	hh := w.highs[0]
	ll := w.lows[0]
	for i := 1; i < w.period; i++ {
		if w.highs[i] > hh {
			hh = w.highs[i]
		}
		if w.lows[i] < ll {
			ll = w.lows[i]
		}
	}

	if hh-ll == 0 {
		return 0
	}
	return ((hh - candle.Close) / (hh - ll)) * -100
}

func (w *WilliamsR) Reset() {
	w.highs = make([]float64, w.period)
	w.lows = make([]float64, w.period)
	w.index = 0
	w.count = 0
}

func (w *WilliamsR) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Williams %R",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(w.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "%R",
				Color:  "#F44336",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: -100, Max: 0},
				Levels: []indicators.Level{
					{Value: -20, Label: "Overbought", Color: "#EF5350"},
					{Value: -80, Label: "Oversold", Color: "#66BB6A"},
				},
			},
		},
	}
}
