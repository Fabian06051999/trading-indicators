package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// StochRSI implements the Stochastic RSI.
type StochRSI struct {
	out       []float64
	rsiPeriod int
	stochK    int
	stochD    int
	rsi       *RSI
	rsiBuffer []float64
	kBuffer   []float64
	rsiIndex  int
	kIndex    int
	rsiCount  int
	kCount    int
}

func NewStochRSI(rsiPeriod, stochK, stochD int) *StochRSI {
	return &StochRSI{
		rsiPeriod: rsiPeriod,
		stochK:    stochK,
		stochD:    stochD,
		rsi:       NewRSI(rsiPeriod),
		rsiBuffer: make([]float64, stochK),
		kBuffer:   make([]float64, stochD),
		out:       make([]float64, 2),
	}
}

func (s *StochRSI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	kValues := make([]float64, len(candles))
	dValues := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		values := s.UpdateAll(c)
		kValues[i] = values[0]
		dValues[i] = values[1]
	}
	return [][]float64{kValues, dValues}
}

func (s *StochRSI) UpdateAll(candle indicators.OHLCV) []float64 {
	rsiVal := s.rsi.UpdateAll(candle)[0]
	if math.IsNaN(rsiVal) {
		s.out[0] = math.NaN()
		s.out[1] = math.NaN()
		return s.out
	}

	s.rsiBuffer[s.rsiIndex] = rsiVal
	s.rsiIndex = (s.rsiIndex + 1) % s.stochK
	s.rsiCount++

	if s.rsiCount < s.stochK {
		s.out[0] = math.NaN()
		s.out[1] = math.NaN()
		return s.out
	}

	// Find min/max RSI in window
	minRSI := s.rsiBuffer[0]
	maxRSI := s.rsiBuffer[0]
	for i := 1; i < s.stochK; i++ {
		if s.rsiBuffer[i] < minRSI {
			minRSI = s.rsiBuffer[i]
		}
		if s.rsiBuffer[i] > maxRSI {
			maxRSI = s.rsiBuffer[i]
		}
	}

	k := 0.0
	if maxRSI-minRSI != 0 {
		k = (rsiVal - minRSI) / (maxRSI - minRSI) * 100
	}

	// %D (SMA of %K)
	s.kBuffer[s.kIndex] = k
	s.kIndex = (s.kIndex + 1) % s.stochD
	s.kCount++

	if s.kCount < s.stochD {
		s.out[0] = k
		s.out[1] = 0
		return s.out
	}

	d := 0.0
	for i := 0; i < s.stochD; i++ {
		d += s.kBuffer[i]
	}
	d /= float64(s.stochD)

	s.out[0] = k
	s.out[1] = d
	return s.out
}

func (s *StochRSI) Reset() {
	s.rsi.Reset()
	s.rsiBuffer = make([]float64, s.stochK)
	s.kBuffer = make([]float64, s.stochD)
	s.rsiIndex = 0
	s.kIndex = 0
	s.rsiCount = 0
	s.kCount = 0
}

func (s *StochRSI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Stochastic RSI",
		Parameters: []indicators.Parameter{
			{Name: "RSI Period", DefaultValue: float64(s.rsiPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Stoch K", DefaultValue: float64(s.stochK), Min: 2, Max: 100, Step: 1},
			{Name: "Stoch D", DefaultValue: float64(s.stochD), Min: 2, Max: 100, Step: 1},
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
			{Name: "%D", Color: "#FF9800", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
