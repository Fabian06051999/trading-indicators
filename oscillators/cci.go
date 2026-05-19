package oscillators

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// CCI implements the Commodity Channel Index.
type CCI struct {
	period int
	buffer []float64
	index  int
	count  int
}

func NewCCI(period int) *CCI {
	period = indicators.ClampMin(period, 2)
	return &CCI{
		period: period,
		buffer: make([]float64, period),
	}
}

func (c *CCI) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		result[i] = c.Update(candle)
	}
	return result
}

func (c *CCI) Update(candle indicators.OHLCV) float64 {
	tp := (candle.High + candle.Low + candle.Close) / 3.0

	c.buffer[c.index] = tp
	c.index = (c.index + 1) % c.period
	if c.count < c.period {
		c.count++
	}
	if c.count < c.period {
		return math.NaN()
	}

	// SMA of typical price
	sum := 0.0
	for i := 0; i < c.period; i++ {
		sum += c.buffer[i]
	}
	sma := sum / float64(c.period)

	// Mean deviation
	md := 0.0
	for i := 0; i < c.period; i++ {
		md += math.Abs(c.buffer[i] - sma)
	}
	md /= float64(c.period)

	if md == 0 {
		return math.NaN()
	}
	return (tp - sma) / (0.015 * md)
}

func (c *CCI) Reset() {
	c.buffer = make([]float64, c.period)
	c.index = 0
	c.count = 0
}

func (c *CCI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Commodity Channel Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(c.period), Min: 2, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "CCI",
				Color: "#009688",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 100, Label: "Overbought", Color: "#EF5350"},
					{Value: -100, Label: "Oversold", Color: "#66BB6A"},
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
