package volume

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// EaseOfMovement implements the Ease of Movement indicator.
type EaseOfMovement struct {
	period    int
	buffer    []float64
	prevHigh  float64
	prevLow   float64
	sum       float64
	index     int
	count     int
}

func NewEaseOfMovement(period int) *EaseOfMovement {
	return &EaseOfMovement{
		period: period,
		buffer: make([]float64, period),
	}
}

func (e *EaseOfMovement) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	e.Reset()

	for i, c := range candles {
		result[i] = e.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (e *EaseOfMovement) UpdateAll(candle indicators.OHLCV) []float64 {
	e.count++

	if e.count == 1 {
		e.prevHigh = candle.High
		e.prevLow = candle.Low
		return []float64{math.NaN()}
	}

	dm := ((candle.High + candle.Low) / 2.0) - ((e.prevHigh + e.prevLow) / 2.0)
	br := 0.0
	hl := candle.High - candle.Low
	if hl != 0 && candle.Volume != 0 {
		br = (candle.Volume / 10000.0) / hl
	}

	e.prevHigh = candle.High
	e.prevLow = candle.Low

	emv := 0.0
	if br != 0 {
		emv = dm / br
	}

	old := e.buffer[e.index]
	e.buffer[e.index] = emv
	e.sum += emv - old
	e.index = (e.index + 1) % e.period

	if e.count <= e.period {
		return []float64{math.NaN()}
	}

	return []float64{e.sum / float64(e.period)}
}

func (e *EaseOfMovement) Reset() {
	e.buffer = make([]float64, e.period)
	e.prevHigh = 0
	e.prevLow = 0
	e.sum = 0
	e.index = 0
	e.count = 0
}

func (e *EaseOfMovement) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Ease of Movement",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(e.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "EMV",
				Color: "#43A047",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
