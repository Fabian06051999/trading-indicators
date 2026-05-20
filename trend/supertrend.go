package trend

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// Supertrend implements the Supertrend indicator.
type Supertrend struct {
	out       []float64
	period    int
	multi     float64
	atrVal    float64
	prevClose float64
	prevATR   float64
	upper     float64
	lower     float64
	trend     int // 1 = up, -1 = down
	prevUpper float64
	prevLower float64
	count     int
	trSum     float64
}

func NewSupertrend(period int, multiplier float64) *Supertrend {
	return &Supertrend{
		period: period,
		multi:  multiplier,
		trend:  1,
		out:    make([]float64, 2),
	}
}

func (s *Supertrend) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	stValues := make([]float64, len(candles))
	trendValues := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		values := s.UpdateAll(c)
		stValues[i] = values[0]
		trendValues[i] = values[1]
	}
	return [][]float64{stValues, trendValues}
}

func (s *Supertrend) UpdateAll(candle indicators.OHLCV) []float64 {
	s.count++

	if s.count == 1 {
		s.prevClose = candle.Close
		s.trSum = candle.High - candle.Low
		s.out[0] = math.NaN()
		s.out[1] = math.NaN()
		return s.out
	}

	// True Range
	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-s.prevClose), math.Abs(candle.Low-s.prevClose)))

	if s.count <= s.period {
		s.trSum += tr
		s.prevClose = candle.Close
		if s.count == s.period {
			s.atrVal = s.trSum / float64(s.period)
		}
		s.out[0] = math.NaN()
		s.out[1] = math.NaN()
		return s.out
	}

	// Wilder's ATR smoothing
	s.atrVal = (s.atrVal*float64(s.period-1) + tr) / float64(s.period)

	hl2 := (candle.High + candle.Low) / 2.0
	upperBand := hl2 + s.multi*s.atrVal
	lowerBand := hl2 - s.multi*s.atrVal

	// Adjust bands
	if lowerBand > s.prevLower || s.prevClose < s.prevLower {
		s.lower = lowerBand
	} else {
		s.lower = s.prevLower
	}

	if upperBand < s.prevUpper || s.prevClose > s.prevUpper {
		s.upper = upperBand
	} else {
		s.upper = s.prevUpper
	}

	// Determine trend
	if s.trend == 1 {
		if candle.Close < s.lower {
			s.trend = -1
		}
	} else {
		if candle.Close > s.upper {
			s.trend = 1
		}
	}

	s.prevUpper = s.upper
	s.prevLower = s.lower
	s.prevClose = candle.Close

	st := s.lower
	if s.trend == -1 {
		st = s.upper
	}
	s.out[0] = st
	s.out[1] = float64(s.trend)
	return s.out
}

func (s *Supertrend) Reset() {
	s.atrVal = 0
	s.prevClose = 0
	s.upper = 0
	s.lower = 0
	s.prevUpper = 0
	s.prevLower = 0
	s.trend = 1
	s.count = 0
	s.trSum = 0
}

func (s *Supertrend) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Supertrend",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(s.period), Min: 2, Max: 100, Step: 1},
			{Name: "Multiplier", DefaultValue: s.multi, Min: 0.5, Max: 10, Step: 0.5},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Supertrend", Color: "#4CAF50", Style: indicators.StyleLine, Width: 2},
			{Name: "Trend Direction", Color: "#9E9E9E", Style: indicators.StyleLine, Width: 0},
		},
	}
}
