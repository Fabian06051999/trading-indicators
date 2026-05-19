package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
)

// Stochastic implements the Stochastic Oscillator (%K and %D).
type Stochastic struct {
	kPeriod  int
	dPeriod  int
	slowing  int
	highs    []float64
	lows     []float64
	rawK     []float64
	index    int
	kIndex   int
	count    int
	kCount   int
}

func NewStochastic(kPeriod, dPeriod, slowing int) *Stochastic {
	return &Stochastic{
		kPeriod: kPeriod,
		dPeriod: dPeriod,
		slowing: slowing,
		highs:   make([]float64, kPeriod),
		lows:    make([]float64, kPeriod),
		rawK:    make([]float64, dPeriod),
	}
}

func (s *Stochastic) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	kValues := make([]float64, len(candles))
	dValues := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		k, d := s.updateBoth(c)
		kValues[i] = k
		dValues[i] = d
	}
	return [][]float64{kValues, dValues}
}

func (s *Stochastic) UpdateAll(candle indicators.OHLCV) []float64 {
	k, d := s.updateBoth(candle)
	return []float64{k, d}
}

func (s *Stochastic) updateBoth(candle indicators.OHLCV) (float64, float64) {
	s.count++
	s.highs[s.index] = candle.High
	s.lows[s.index] = candle.Low
	s.index = (s.index + 1) % s.kPeriod

	if s.count < s.kPeriod {
		return 0, 0
	}

	// Find highest high and lowest low
	hh := s.highs[0]
	ll := s.lows[0]
	for i := 1; i < s.kPeriod; i++ {
		if s.highs[i] > hh {
			hh = s.highs[i]
		}
		if s.lows[i] < ll {
			ll = s.lows[i]
		}
	}

	// %K
	k := 0.0
	if hh-ll != 0 {
		k = (candle.Close - ll) / (hh - ll) * 100
	}

	// %D (SMA of %K)
	s.rawK[s.kIndex] = k
	s.kIndex = (s.kIndex + 1) % s.dPeriod
	s.kCount++

	if s.kCount < s.dPeriod {
		return k, 0
	}

	d := 0.0
	for i := 0; i < s.dPeriod; i++ {
		d += s.rawK[i]
	}
	d /= float64(s.dPeriod)

	return k, d
}

func (s *Stochastic) Reset() {
	s.highs = make([]float64, s.kPeriod)
	s.lows = make([]float64, s.kPeriod)
	s.rawK = make([]float64, s.dPeriod)
	s.index = 0
	s.kIndex = 0
	s.count = 0
	s.kCount = 0
}

func (s *Stochastic) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Stochastic Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "K Period", DefaultValue: float64(s.kPeriod), Min: 1, Max: 100, Step: 1},
			{Name: "D Period", DefaultValue: float64(s.dPeriod), Min: 1, Max: 100, Step: 1},
			{Name: "Slowing", DefaultValue: float64(s.slowing), Min: 1, Max: 10, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "%K",
				Color:  "#2196F3",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 80, Label: "Overbought", Color: "#EF5350"},
					{Value: 20, Label: "Oversold", Color: "#66BB6A"},
				},
			},
			{
				Name:  "%D",
				Color: "#FF9800",
				Style: indicators.StyleDashed,
				Width: 1,
			},
		},
	}
}
