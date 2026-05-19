package moving_averages

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// LSMA implements the Least Squares Moving Average (Linear Regression).
type LSMA struct {
	period int
	buffer []float64
	index  int
	count  int
}

func NewLSMA(period int) *LSMA {
	period = indicators.ClampMin(period, 2)
	return &LSMA{
		period: period,
		buffer: make([]float64, period),
	}
}

func (l *LSMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	l.Reset()

	for i, c := range candles {
		result[i] = l.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (l *LSMA) UpdateAll(candle indicators.OHLCV) []float64 {
	l.buffer[l.index] = candle.Close
	l.index = (l.index + 1) % l.period
	if l.count < l.period {
		l.count++
	}
	if l.count < l.period {
		return []float64{math.NaN()}
	}

	// Linear regression using least squares
	n := float64(l.period)
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i := 0; i < l.period; i++ {
		idx := (l.index + i) % l.period
		x := float64(i + 1)
		y := l.buffer[idx]
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return []float64{math.NaN()}
	}

	slope := (n*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / n

	return []float64{intercept + slope*n}
}

func (l *LSMA) Reset() {
	l.buffer = make([]float64, l.period)
	l.index = 0
	l.count = 0
}

func (l *LSMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Least Squares Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(l.period), Min: 2, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "LSMA", Color: "#CDDC39", Style: indicators.StyleLine, Width: 2},
		},
	}
}
