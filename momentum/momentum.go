package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// Momentum implements the basic Momentum indicator (price - price[n]).
type Momentum struct {
	out    []float64
	period int
	buffer []float64
	index  int
	count  int
}

func NewMomentum(period int) *Momentum {
	period = indicators.ClampMin(period, 1)
	return &Momentum{
		period: period,
		buffer: make([]float64, period+1),
		out:    make([]float64, 1),
	}
}

func (m *Momentum) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	m.Reset()

	for i, c := range candles {
		result[i] = m.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (m *Momentum) UpdateAll(candle indicators.OHLCV) []float64 {
	m.buffer[m.index] = candle.Close
	m.count++

	if m.count <= m.period {
		m.index = (m.index + 1) % (m.period + 1)
		m.out[0] = math.NaN()
		return m.out
	}

	pastIdx := (m.index + 1) % (m.period + 1)
	result := candle.Close - m.buffer[pastIdx]
	m.index = (m.index + 1) % (m.period + 1)
	m.out[0] = result
	return m.out
}

func (m *Momentum) Reset() {
	m.buffer = make([]float64, m.period+1)
	m.index = 0
	m.count = 0
}

func (m *Momentum) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Momentum",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(m.period), Min: 1, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Momentum",
				Color: "#D84315",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
