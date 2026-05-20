package momentum

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/volatility"
)

// SqueezeMomentum implements the Squeeze Momentum Indicator (TTM Squeeze).
type SqueezeMomentum struct {
	out       []float64
	bbPeriod  int
	bbStdDev  float64
	kcPeriod  int
	kcMulti   float64
	bb        *volatility.BollingerBands
	kc        *volatility.KeltnerChannel
	linRegBuf []float64
	index     int
	count     int
}

func NewSqueezeMomentum(bbPeriod int, bbStdDev float64, kcPeriod int, kcMulti float64) *SqueezeMomentum {
	return &SqueezeMomentum{
		bbPeriod:  bbPeriod,
		bbStdDev:  bbStdDev,
		kcPeriod:  kcPeriod,
		kcMulti:   kcMulti,
		bb:        volatility.NewBollingerBands(bbPeriod, bbStdDev),
		kc:        volatility.NewKeltnerChannel(kcPeriod, kcPeriod, kcMulti),
		linRegBuf: make([]float64, bbPeriod),
		out:       make([]float64, 2),
	}
}

// CalculateAll returns [momentum, squeeze_on (1=squeeze, 0=no squeeze)]
func (s *SqueezeMomentum) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	momValues := make([]float64, len(candles))
	sqzValues := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		values := s.UpdateAll(c)
		momValues[i] = values[0]
		sqzValues[i] = values[1]
	}
	return [][]float64{momValues, sqzValues}
}

func (s *SqueezeMomentum) UpdateAll(candle indicators.OHLCV) []float64 {
	bbVals := s.bb.UpdateAll(candle)
	kcVals := s.kc.UpdateAll(candle)

	s.count++

	// Store close-midline value for linear regression
	mid := (bbVals[0] + bbVals[2]) / 2.0
	val := candle.Close - mid
	s.linRegBuf[s.index] = val
	s.index = (s.index + 1) % s.bbPeriod

	if bbVals[0] == 0 || kcVals[0] == 0 {
		s.out[0] = math.NaN()
		s.out[1] = math.NaN()
		return s.out
	}

	// Squeeze: BB inside KC
	sqz := 0.0
	if bbVals[2] > kcVals[2] && bbVals[0] < kcVals[0] {
		sqz = 1.0
	}

	// Momentum: linear regression of (close - midline(BB))
	if s.count < s.bbPeriod {
		s.out[0] = math.NaN()
		s.out[1] = sqz
		return s.out
	}

	// Simple linear regression value
	n := float64(s.bbPeriod)
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	for i := 0; i < s.bbPeriod; i++ {
		idx := (s.index + i) % s.bbPeriod
		x := float64(i + 1)
		y := s.linRegBuf[idx]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denom := n*sumX2 - sumX*sumX
	mom := 0.0
	if denom != 0 {
		slope := (n*sumXY - sumX*sumY) / denom
		intercept := (sumY - slope*sumX) / n
		mom = intercept + slope*n
	}

	s.out[0] = mom
	s.out[1] = sqz
	return s.out
}

func (s *SqueezeMomentum) Reset() {
	s.bb.Reset()
	s.kc.Reset()
	s.linRegBuf = make([]float64, s.bbPeriod)
	s.index = 0
	s.count = 0
}

func (s *SqueezeMomentum) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Squeeze Momentum",
		Parameters: []indicators.Parameter{
			{Name: "BB Period", DefaultValue: float64(s.bbPeriod), Min: 5, Max: 50, Step: 1},
			{Name: "BB StdDev", DefaultValue: s.bbStdDev, Min: 0.5, Max: 5, Step: 0.5},
			{Name: "KC Period", DefaultValue: float64(s.kcPeriod), Min: 5, Max: 50, Step: 1},
			{Name: "KC Multiplier", DefaultValue: s.kcMulti, Min: 0.5, Max: 5, Step: 0.5},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Momentum",
				Color: "#00BFA5",
				Style: indicators.StyleHistogram,
				Width: 1,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
			{Name: "Squeeze", Color: "#FF1744", Style: indicators.StyleDots, Width: 3},
		},
	}
}
