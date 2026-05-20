package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// PivotPoints implements Standard (Floor) Pivot Points.
type PivotPoints struct {
	out       []float64
	prevHigh  float64
	prevLow   float64
	prevClose float64
	count     int
}

func NewPivotPoints() *PivotPoints {
	return &PivotPoints{out: make([]float64, 7)}
}

// CalculateAll returns [PP, R1, R2, R3, S1, S2, S3]
func (p *PivotPoints) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	pp := make([]float64, len(candles))
	r1 := make([]float64, len(candles))
	r2 := make([]float64, len(candles))
	r3 := make([]float64, len(candles))
	s1 := make([]float64, len(candles))
	s2 := make([]float64, len(candles))
	s3 := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		values := p.UpdateAll(c)
		pp[i] = values[0]
		r1[i] = values[1]
		r2[i] = values[2]
		r3[i] = values[3]
		s1[i] = values[4]
		s2[i] = values[5]
		s3[i] = values[6]
	}
	return [][]float64{pp, r1, r2, r3, s1, s2, s3}
}

func (p *PivotPoints) UpdateAll(candle indicators.OHLCV) []float64 {
	p.count++

	if p.count == 1 {
		p.prevHigh = candle.High
		p.prevLow = candle.Low
		p.prevClose = candle.Close
		p.out[0] = math.NaN()
		p.out[1] = math.NaN()
		p.out[2] = math.NaN()
		p.out[3] = math.NaN()
		p.out[4] = math.NaN()
		p.out[5] = math.NaN()
		p.out[6] = math.NaN()
		return p.out
	}

	// Calculate pivots based on previous candle
	pivot := (p.prevHigh + p.prevLow + p.prevClose) / 3.0
	r1Val := 2*pivot - p.prevLow
	s1Val := 2*pivot - p.prevHigh
	r2Val := pivot + (p.prevHigh - p.prevLow)
	s2Val := pivot - (p.prevHigh - p.prevLow)
	r3Val := p.prevHigh + 2*(pivot-p.prevLow)
	s3Val := p.prevLow - 2*(p.prevHigh-pivot)

	p.prevHigh = candle.High
	p.prevLow = candle.Low
	p.prevClose = candle.Close

	p.out[0] = pivot
	p.out[1] = r1Val
	p.out[2] = r2Val
	p.out[3] = r3Val
	p.out[4] = s1Val
	p.out[5] = s2Val
	p.out[6] = s3Val
	return p.out
}

func (p *PivotPoints) Reset() {
	p.prevHigh = 0
	p.prevLow = 0
	p.prevClose = 0
	p.count = 0
}

func (p *PivotPoints) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name:       "Pivot Points (Standard)",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "PP", Color: "#FFB300", Style: indicators.StyleDashed, Width: 2},
			{Name: "R1", Color: "#EF5350", Style: indicators.StyleDashed, Width: 1},
			{Name: "R2", Color: "#E53935", Style: indicators.StyleDashed, Width: 1},
			{Name: "R3", Color: "#C62828", Style: indicators.StyleDashed, Width: 1},
			{Name: "S1", Color: "#66BB6A", Style: indicators.StyleDashed, Width: 1},
			{Name: "S2", Color: "#43A047", Style: indicators.StyleDashed, Width: 1},
			{Name: "S3", Color: "#2E7D32", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
