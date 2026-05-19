package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// StdDev implements the Standard Deviation indicator.
type StdDev struct {
	period int
	buffer []float64
	sum    float64
	index  int
	count  int
}

func NewStdDev(period int) *StdDev {
	return &StdDev{
		period: period,
		buffer: make([]float64, period),
	}
}

func (s *StdDev) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		result[i] = s.Update(c)
	}
	return result
}

func (s *StdDev) Update(candle indicators.OHLCV) float64 {
	old := s.buffer[s.index]
	s.buffer[s.index] = candle.Close
	s.sum += candle.Close - old
	s.index = (s.index + 1) % s.period
	if s.count < s.period {
		s.count++
	}
	if s.count < s.period {
		return 0
	}

	mean := s.sum / float64(s.period)
	variance := 0.0
	for i := 0; i < s.period; i++ {
		diff := s.buffer[i] - mean
		variance += diff * diff
	}
	return math.Sqrt(variance / float64(s.period))
}

func (s *StdDev) Reset() {
	s.buffer = make([]float64, s.period)
	s.sum = 0
	s.index = 0
	s.count = 0
}

func (s *StdDev) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Standard Deviation",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(s.period), Min: 2, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "StdDev", Color: "#FF7043", Style: indicators.StyleLine, Width: 2},
		},
	}
}
