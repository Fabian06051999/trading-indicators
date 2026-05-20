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
	out    []float64
}

func NewCCI(period int) *CCI {
	period = indicators.ClampMin(period, 2)
	return &CCI{
		period: period,
		buffer: make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (c *CCI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		c.update(candle)
		result[i] = c.out[0]
	}
	return [][]float64{result}
}

func (c *CCI) UpdateAll(candle indicators.OHLCV) []float64 {
	c.update(candle)
	return c.out
}

func (c *CCI) update(candle indicators.OHLCV) {
	tp := (candle.High + candle.Low + candle.Close) / 3.0

	c.buffer[c.index] = tp
	c.index = (c.index + 1) % c.period
	if c.count < c.period {
		c.count++
	}
	if c.count < c.period {
		c.out[0] = math.NaN()
		return
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
		c.out[0] = 0
		return
	}
	c.out[0] = (tp - sma) / (0.015 * md)
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
