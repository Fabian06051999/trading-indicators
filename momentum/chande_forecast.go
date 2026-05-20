package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// ChandeForecast implements the Chande Forecast Oscillator.
// CFO = ((Close - LinearRegression(Close, period)) / Close) * 100
type ChandeForecast struct {
	out    []float64
	period int
	buffer []float64
	index  int
	count  int
}

func NewChandeForecast(period int) *ChandeForecast {
	return &ChandeForecast{
		period: period,
		buffer: make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (cf *ChandeForecast) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	cf.Reset()

	for i, c := range candles {
		result[i] = cf.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (cf *ChandeForecast) UpdateAll(candle indicators.OHLCV) []float64 {
	cf.buffer[cf.index] = candle.Close
	cf.index = (cf.index + 1) % cf.period
	if cf.count < cf.period {
		cf.count++
	}
	if cf.count < cf.period {
		cf.out[0] = math.NaN()
		return cf.out
	}

	// Linear regression forecast
	n := float64(cf.period)
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i := 0; i < cf.period; i++ {
		idx := (cf.index + i) % cf.period
		x := float64(i + 1)
		y := cf.buffer[idx]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denom := n*sumX2 - sumX*sumX
	if denom == 0 || candle.Close == 0 {
		cf.out[0] = math.NaN()
		return cf.out
	}

	slope := (n*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / n
	forecast := intercept + slope*(n+1)

	cf.out[0] = ((candle.Close - forecast) / candle.Close) * 100
	return cf.out
}

func (cf *ChandeForecast) Reset() {
	cf.buffer = make([]float64, cf.period)
	cf.index = 0
	cf.count = 0
}

func (cf *ChandeForecast) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Chande Forecast Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(cf.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "CFO",
				Color: "#6A1B9A",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
