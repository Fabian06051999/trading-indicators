package trend

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// Aroon implements the Aroon Up/Down indicator.
type Aroon struct {
	out    []float64
	period int
	highs  []float64
	lows   []float64
	index  int
	count  int
}

func NewAroon(period int) *Aroon {
	return &Aroon{
		period: period,
		highs:  make([]float64, period+1),
		lows:   make([]float64, period+1),
		out:    make([]float64, 2),
	}
}

func (a *Aroon) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	upValues := make([]float64, len(candles))
	downValues := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		values := a.UpdateAll(c)
		upValues[i] = values[0]
		downValues[i] = values[1]
	}
	return [][]float64{upValues, downValues}
}

func (a *Aroon) UpdateAll(candle indicators.OHLCV) []float64 {
	a.highs[a.index] = candle.High
	a.lows[a.index] = candle.Low
	a.count++
	a.index = (a.index + 1) % (a.period + 1)

	if a.count <= a.period {
		a.out[0] = math.NaN()
		a.out[1] = math.NaN()
		return a.out
	}

	// Find index of highest high and lowest low
	highIdx := 0
	lowIdx := 0
	hh := -1e308
	ll := 1e308

	for i := 0; i <= a.period; i++ {
		idx := (a.index + i) % (a.period + 1)
		if a.highs[idx] >= hh {
			hh = a.highs[idx]
			highIdx = i
		}
		if a.lows[idx] <= ll {
			ll = a.lows[idx]
			lowIdx = i
		}
	}

	p := float64(a.period)
	aroonUp := (float64(highIdx) / p) * 100
	aroonDown := (float64(lowIdx) / p) * 100

	a.out[0] = aroonUp
	a.out[1] = aroonDown
	return a.out
}

func (a *Aroon) Reset() {
	a.highs = make([]float64, a.period+1)
	a.lows = make([]float64, a.period+1)
	a.index = 0
	a.count = 0
}

func (a *Aroon) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Aroon",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(a.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "Aroon Up",
				Color:  "#4CAF50",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
			},
			{
				Name:  "Aroon Down",
				Color: "#F44336",
				Style: indicators.StyleLine,
				Width: 2,
			},
		},
	}
}
