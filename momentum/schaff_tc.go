package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// SchaffTrendCycle implements the Schaff Trend Cycle.
type SchaffTrendCycle struct {
	macdFast   int
	macdSlow   int
	cyclePeriod int
	fastEMA    *moving_averages.EMA
	slowEMA    *moving_averages.EMA
	macdBuf    []float64
	stoch1Buf  []float64
	macdIdx    int
	stoch1Idx  int
	macdCount  int
	stoch1Count int
	pf         float64
	pff        float64
}

func NewSchaffTrendCycle(macdFast, macdSlow, cyclePeriod int) *SchaffTrendCycle {
	return &SchaffTrendCycle{
		macdFast:    macdFast,
		macdSlow:    macdSlow,
		cyclePeriod: cyclePeriod,
		fastEMA:     moving_averages.NewEMA(macdFast),
		slowEMA:     moving_averages.NewEMA(macdSlow),
		macdBuf:     make([]float64, cyclePeriod),
		stoch1Buf:   make([]float64, cyclePeriod),
	}
}

func (s *SchaffTrendCycle) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	s.Reset()

	for i, c := range candles {
		result[i] = s.Update(c)
	}
	return result
}

func (s *SchaffTrendCycle) Update(candle indicators.OHLCV) float64 {
	fast := s.fastEMA.Update(candle)
	slow := s.slowEMA.Update(candle)

	if fast == 0 || slow == 0 {
		return 0
	}

	macdVal := fast - slow

	// First Stochastic of MACD
	s.macdBuf[s.macdIdx] = macdVal
	s.macdIdx = (s.macdIdx + 1) % s.cyclePeriod
	s.macdCount++

	if s.macdCount < s.cyclePeriod {
		return 0
	}

	ll := s.macdBuf[0]
	hh := s.macdBuf[0]
	for i := 1; i < s.cyclePeriod; i++ {
		if s.macdBuf[i] < ll {
			ll = s.macdBuf[i]
		}
		if s.macdBuf[i] > hh {
			hh = s.macdBuf[i]
		}
	}

	stoch1 := 0.0
	if hh-ll != 0 {
		stoch1 = ((macdVal - ll) / (hh - ll)) * 100
	}

	// Smoothed first stochastic
	s.pf = s.pf + 0.5*(stoch1-s.pf)

	// Second Stochastic
	s.stoch1Buf[s.stoch1Idx] = s.pf
	s.stoch1Idx = (s.stoch1Idx + 1) % s.cyclePeriod
	s.stoch1Count++

	if s.stoch1Count < s.cyclePeriod {
		return 0
	}

	ll2 := s.stoch1Buf[0]
	hh2 := s.stoch1Buf[0]
	for i := 1; i < s.cyclePeriod; i++ {
		if s.stoch1Buf[i] < ll2 {
			ll2 = s.stoch1Buf[i]
		}
		if s.stoch1Buf[i] > hh2 {
			hh2 = s.stoch1Buf[i]
		}
	}

	stoch2 := 0.0
	if hh2-ll2 != 0 {
		stoch2 = ((s.pf - ll2) / (hh2 - ll2)) * 100
	}

	s.pff = s.pff + 0.5*(stoch2-s.pff)
	return s.pff
}

func (s *SchaffTrendCycle) Reset() {
	s.fastEMA.Reset()
	s.slowEMA.Reset()
	s.macdBuf = make([]float64, s.cyclePeriod)
	s.stoch1Buf = make([]float64, s.cyclePeriod)
	s.macdIdx = 0
	s.stoch1Idx = 0
	s.macdCount = 0
	s.stoch1Count = 0
	s.pf = 0
	s.pff = 0
}

func (s *SchaffTrendCycle) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Schaff Trend Cycle",
		Parameters: []indicators.Parameter{
			{Name: "MACD Fast", DefaultValue: float64(s.macdFast), Min: 2, Max: 50, Step: 1},
			{Name: "MACD Slow", DefaultValue: float64(s.macdSlow), Min: 5, Max: 100, Step: 1},
			{Name: "Cycle Period", DefaultValue: float64(s.cyclePeriod), Min: 2, Max: 50, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "STC",
				Color:  "#1A237E",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 75, Label: "Overbought", Color: "#EF5350"},
					{Value: 25, Label: "Oversold", Color: "#66BB6A"},
				},
			},
		},
	}
}
