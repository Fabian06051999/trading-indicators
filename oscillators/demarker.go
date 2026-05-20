package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// DeMarker implements the DeMarker oscillator.
type DeMarker struct {
	period   int
	deMax    []float64
	deMin    []float64
	prevHigh float64
	prevLow  float64
	index    int
	count    int
	out      []float64
}

func NewDeMarker(period int) *DeMarker {
	period = indicators.ClampMin(period, 2)
	return &DeMarker{
		period: period,
		deMax:  make([]float64, period),
		deMin:  make([]float64, period),
		out:    make([]float64, 1),
	}
}

func (d *DeMarker) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		d.update(c)
		result[i] = d.out[0]
	}
	return [][]float64{result}
}

func (d *DeMarker) UpdateAll(candle indicators.OHLCV) []float64 {
	d.update(candle)
	return d.out
}

func (d *DeMarker) update(candle indicators.OHLCV) {
	d.count++

	if d.count == 1 {
		d.prevHigh = candle.High
		d.prevLow = candle.Low
		d.out[0] = math.NaN()
		return
	}

	deMax := 0.0
	if candle.High > d.prevHigh {
		deMax = candle.High - d.prevHigh
	}

	deMin := 0.0
	if candle.Low < d.prevLow {
		deMin = d.prevLow - candle.Low
	}

	d.prevHigh = candle.High
	d.prevLow = candle.Low

	d.deMax[d.index] = deMax
	d.deMin[d.index] = deMin
	d.index = (d.index + 1) % d.period

	if d.count <= d.period {
		d.out[0] = math.NaN()
		return
	}

	sumMax := 0.0
	sumMin := 0.0
	for i := 0; i < d.period; i++ {
		sumMax += d.deMax[i]
		sumMin += d.deMin[i]
	}

	if sumMax+sumMin == 0 {
		d.out[0] = 0.5
		return
	}
	d.out[0] = sumMax / (sumMax + sumMin)
}

func (d *DeMarker) Reset() {
	d.deMax = make([]float64, d.period)
	d.deMin = make([]float64, d.period)
	d.prevHigh = 0
	d.prevLow = 0
	d.index = 0
	d.count = 0
}

func (d *DeMarker) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "DeMarker",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(d.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "DeM",
				Color:  "#5E35B1",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 1},
				Levels: []indicators.Level{
					{Value: 0.7, Label: "Overbought", Color: "#EF5350"},
					{Value: 0.3, Label: "Oversold", Color: "#66BB6A"},
				},
			},
		},
	}
}
