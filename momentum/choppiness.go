package momentum

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// ChoppinessIndex measures if the market is choppy (range-bound) or trending.
type ChoppinessIndex struct {
	out       []float64
	period    int
	trBuffer  []float64
	highs     []float64
	lows      []float64
	prevClose float64
	index     int
	count     int
}

func NewChoppinessIndex(period int) *ChoppinessIndex {
	return &ChoppinessIndex{
		period:   period,
		trBuffer: make([]float64, period),
		highs:    make([]float64, period),
		lows:     make([]float64, period),
		out:      make([]float64, 1),
	}
}

func (ci *ChoppinessIndex) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	ci.Reset()

	for i, c := range candles {
		result[i] = ci.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (ci *ChoppinessIndex) UpdateAll(candle indicators.OHLCV) []float64 {
	ci.count++

	if ci.count == 1 {
		ci.prevClose = candle.Close
		ci.highs[ci.index] = candle.High
		ci.lows[ci.index] = candle.Low
		ci.trBuffer[ci.index] = candle.High - candle.Low
		ci.index = (ci.index + 1) % ci.period
		ci.out[0] = math.NaN()
		return ci.out
	}

	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-ci.prevClose), math.Abs(candle.Low-ci.prevClose)))
	ci.prevClose = candle.Close

	ci.trBuffer[ci.index] = tr
	ci.highs[ci.index] = candle.High
	ci.lows[ci.index] = candle.Low
	ci.index = (ci.index + 1) % ci.period

	if ci.count < ci.period {
		ci.out[0] = math.NaN()
		return ci.out
	}

	// Sum of TR
	sumTR := 0.0
	for i := 0; i < ci.period; i++ {
		sumTR += ci.trBuffer[i]
	}

	// Highest high - lowest low
	hh := ci.highs[0]
	ll := ci.lows[0]
	for i := 1; i < ci.period; i++ {
		if ci.highs[i] > hh {
			hh = ci.highs[i]
		}
		if ci.lows[i] < ll {
			ll = ci.lows[i]
		}
	}

	hl := hh - ll
	if hl == 0 {
		ci.out[0] = math.NaN()
		return ci.out
	}

	ci.out[0] = 100 * math.Log10(sumTR/hl) / math.Log10(float64(ci.period))
	return ci.out
}

func (ci *ChoppinessIndex) Reset() {
	ci.trBuffer = make([]float64, ci.period)
	ci.highs = make([]float64, ci.period)
	ci.lows = make([]float64, ci.period)
	ci.prevClose = 0
	ci.index = 0
	ci.count = 0
}

func (ci *ChoppinessIndex) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Choppiness Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(ci.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "CHOP",
				Color:  "#6D4C41",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 61.8, Label: "Choppy", Color: "#EF5350"},
					{Value: 38.2, Label: "Trending", Color: "#66BB6A"},
				},
			},
		},
	}
}
