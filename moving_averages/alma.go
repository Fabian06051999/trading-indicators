package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// ALMA implements the Arnaud Legoux Moving Average.
type ALMA struct {
	period int
	offset float64
	sigma  float64
	weights []float64
	buffer  []float64
	index   int
	count   int
}

func NewALMA(period int, offset, sigma float64) *ALMA {
	a := &ALMA{
		period: period,
		offset: offset,
		sigma:  sigma,
		buffer: make([]float64, period),
	}
	a.computeWeights()
	return a
}

func (a *ALMA) computeWeights() {
	a.weights = make([]float64, a.period)
	m := a.offset * float64(a.period-1)
	s := float64(a.period) / a.sigma

	sum := 0.0
	for i := 0; i < a.period; i++ {
		w := math.Exp(-((float64(i) - m) * (float64(i) - m)) / (2 * s * s))
		a.weights[i] = w
		sum += w
	}
	for i := range a.weights {
		a.weights[i] /= sum
	}
}

func (a *ALMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		result[i] = a.Update(c)
	}
	return result
}

func (a *ALMA) Update(candle indicators.OHLCV) float64 {
	a.buffer[a.index] = candle.Close
	a.index = (a.index + 1) % a.period
	if a.count < a.period {
		a.count++
	}
	if a.count < a.period {
		return 0
	}

	sum := 0.0
	for i := 0; i < a.period; i++ {
		idx := (a.index + i) % a.period
		sum += a.buffer[idx] * a.weights[i]
	}
	return sum
}

func (a *ALMA) Reset() {
	a.buffer = make([]float64, a.period)
	a.index = 0
	a.count = 0
}

func (a *ALMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Arnaud Legoux Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(a.period), Min: 2, Max: 500, Step: 1},
			{Name: "Offset", DefaultValue: a.offset, Min: 0, Max: 1, Step: 0.05},
			{Name: "Sigma", DefaultValue: a.sigma, Min: 1, Max: 20, Step: 0.5},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "ALMA", Color: "#8BC34A", Style: indicators.StyleLine, Width: 2},
		},
	}
}
