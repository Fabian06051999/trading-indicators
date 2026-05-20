package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// SMA implements the Simple Moving Average.
type SMA struct {
	period int
	buffer []float64
	sum    float64
	index  int
	count  int
	out    []float64
}

func NewSMA(period int) *SMA {
	period = indicators.ClampMin(period, 1)
	return &SMA{
		period: period,
		buffer: make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (s *SMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		s.update(c)
		result[i] = s.out[0]
	}
	return [][]float64{result}
}

func (s *SMA) UpdateAll(candle indicators.OHLCV) []float64 {
	s.update(candle)
	return s.out
}

func (s *SMA) update(candle indicators.OHLCV) {
	old := s.buffer[s.index]
	s.buffer[s.index] = candle.Close
	s.sum += candle.Close - old
	s.index = (s.index + 1) % s.period
	if s.count < s.period {
		s.count++
	}
	if s.count < s.period {
		s.out[0] = math.NaN()
		return
	}
	s.out[0] = s.sum / float64(s.period)
}

func (s *SMA) Reset() {
	s.buffer = make([]float64, s.period)
	s.sum = 0
	s.index = 0
	s.count = 0
}

func (s *SMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Simple Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(s.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "SMA", Color: "#2196F3", Style: indicators.StyleLine, Width: 2},
		},
	}
}
