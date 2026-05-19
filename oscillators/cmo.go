package oscillators

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// CMO implements the Chande Momentum Oscillator.
type CMO struct {
	period    int
	gains     []float64
	losses    []float64
	prevClose float64
	index     int
	count     int
}

func NewCMO(period int) *CMO {
	period = indicators.ClampMin(period, 2)
	return &CMO{
		period: period,
		gains:  make([]float64, period),
		losses: make([]float64, period),
	}
}

func (c *CMO) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		result[i] = c.Update(candle)
	}
	return result
}

func (c *CMO) Update(candle indicators.OHLCV) float64 {
	c.count++

	if c.count == 1 {
		c.prevClose = candle.Close
		return math.NaN()
	}

	change := candle.Close - c.prevClose
	c.prevClose = candle.Close

	c.gains[c.index] = 0
	c.losses[c.index] = 0
	if change > 0 {
		c.gains[c.index] = change
	} else {
		c.losses[c.index] = -change
	}
	c.index = (c.index + 1) % c.period

	if c.count <= c.period {
		return math.NaN()
	}

	sumGains := 0.0
	sumLosses := 0.0
	for i := 0; i < c.period; i++ {
		sumGains += c.gains[i]
		sumLosses += c.losses[i]
	}

	if sumGains+sumLosses == 0 {
		return math.NaN()
	}
	return ((sumGains - sumLosses) / (sumGains + sumLosses)) * 100
}

func (c *CMO) Reset() {
	c.gains = make([]float64, c.period)
	c.losses = make([]float64, c.period)
	c.prevClose = 0
	c.index = 0
	c.count = 0
}

func (c *CMO) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Chande Momentum Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(c.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "CMO",
				Color:  "#FF7043",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: -100, Max: 100},
				Levels: []indicators.Level{
					{Value: 50, Label: "Overbought", Color: "#EF5350"},
					{Value: -50, Label: "Oversold", Color: "#66BB6A"},
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
