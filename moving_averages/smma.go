package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// SMMA implements the Smoothed Moving Average (also known as RMA/Modified MA).
type SMMA struct {
	period int
	value  float64
	count  int
	sum    float64
	out    []float64
}

func NewSMMA(period int) *SMMA {
	period = indicators.ClampMin(period, 1)
	return &SMMA{
		period: period,
		out:    make([]float64, 1),
	}
}

func (s *SMMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		s.update(c)
		result[i] = s.out[0]
	}
	return [][]float64{result}
}

func (s *SMMA) UpdateAll(candle indicators.OHLCV) []float64 {
	s.update(candle)
	return s.out
}

func (s *SMMA) update(candle indicators.OHLCV) {
	s.count++
	if s.count <= s.period {
		s.sum += candle.Close
		if s.count == s.period {
			s.value = s.sum / float64(s.period)
			s.out[0] = s.value
			return
		}
		s.out[0] = math.NaN()
		return
	}
	s.value = (s.value*float64(s.period-1) + candle.Close) / float64(s.period)
	s.out[0] = s.value
}

func (s *SMMA) Reset() {
	s.value = 0
	s.count = 0
	s.sum = 0
}

func (s *SMMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Smoothed Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(s.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "SMMA", Color: "#795548", Style: indicators.StyleLine, Width: 2},
		},
	}
}
