package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// WilliamsR implements Williams %R.
type WilliamsR struct {
	period int
	highs  []float64
	lows   []float64
	index  int
	count  int
	out    []float64
}

func NewWilliamsR(period int) *WilliamsR {
	period = indicators.ClampMin(period, 2)
	return &WilliamsR{
		period: period,
		highs:  make([]float64, period),
		lows:   make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (w *WilliamsR) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		w.update(c)
		result[i] = w.out[0]
	}
	return [][]float64{result}
}

func (w *WilliamsR) UpdateAll(candle indicators.OHLCV) []float64 {
	w.update(candle)
	return w.out
}

func (w *WilliamsR) update(candle indicators.OHLCV) {
	w.highs[w.index] = candle.High
	w.lows[w.index] = candle.Low
	w.index = (w.index + 1) % w.period
	if w.count < w.period {
		w.count++
	}
	if w.count < w.period {
		w.out[0] = math.NaN()
		return
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
		w.out[0] = math.NaN()
		return
	}
	w.out[0] = ((hh - candle.Close) / (hh - ll)) * -100
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
