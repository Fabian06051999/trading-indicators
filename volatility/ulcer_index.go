package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// UlcerIndex implements the Ulcer Index (downside volatility).
type UlcerIndex struct {
	period int
	buffer []float64
	index  int
	count  int
	maxVal float64
}

func NewUlcerIndex(period int) *UlcerIndex {
	return &UlcerIndex{
		period: period,
		buffer: make([]float64, period),
	}
}

func (u *UlcerIndex) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	u.Reset()

	for i, c := range candles {
		result[i] = u.Update(c)
	}
	return result
}

func (u *UlcerIndex) Update(candle indicators.OHLCV) float64 {
	u.count++

	if candle.Close > u.maxVal {
		u.maxVal = candle.Close
	}

	pctDrawdown := 0.0
	if u.maxVal != 0 {
		pctDrawdown = ((candle.Close - u.maxVal) / u.maxVal) * 100
	}

	u.buffer[u.index] = pctDrawdown * pctDrawdown
	u.index = (u.index + 1) % u.period

	if u.count < u.period {
		return 0
	}

	sum := 0.0
	for i := 0; i < u.period; i++ {
		sum += u.buffer[i]
	}
	return math.Sqrt(sum / float64(u.period))
}

func (u *UlcerIndex) Reset() {
	u.buffer = make([]float64, u.period)
	u.index = 0
	u.count = 0
	u.maxVal = 0
}

func (u *UlcerIndex) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Ulcer Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(u.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "UI", Color: "#C62828", Style: indicators.StyleLine, Width: 2},
		},
	}
}
