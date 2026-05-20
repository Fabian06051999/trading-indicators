package trend

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// Vortex implements the Vortex Indicator (VI+ and VI-).
type Vortex struct {
	out       []float64
	period    int
	vmPlus    []float64
	vmMinus   []float64
	trValues  []float64
	prevHigh  float64
	prevLow   float64
	prevClose float64
	index     int
	count     int
}

func NewVortex(period int) *Vortex {
	return &Vortex{
		period:   period,
		vmPlus:   make([]float64, period),
		vmMinus:  make([]float64, period),
		trValues: make([]float64, period),
		out:      make([]float64, 2),
	}
}

func (v *Vortex) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	viPlus := make([]float64, len(candles))
	viMinus := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		values := v.UpdateAll(c)
		viPlus[i] = values[0]
		viMinus[i] = values[1]
	}
	return [][]float64{viPlus, viMinus}
}

func (v *Vortex) UpdateAll(candle indicators.OHLCV) []float64 {
	v.count++

	if v.count == 1 {
		v.prevHigh = candle.High
		v.prevLow = candle.Low
		v.prevClose = candle.Close
		v.out[0] = math.NaN()
		v.out[1] = math.NaN()
		return v.out
	}

	vmP := math.Abs(candle.High - v.prevLow)
	vmM := math.Abs(candle.Low - v.prevHigh)
	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-v.prevClose), math.Abs(candle.Low-v.prevClose)))

	v.prevHigh = candle.High
	v.prevLow = candle.Low
	v.prevClose = candle.Close

	v.vmPlus[v.index] = vmP
	v.vmMinus[v.index] = vmM
	v.trValues[v.index] = tr
	v.index = (v.index + 1) % v.period

	if v.count <= v.period {
		v.out[0] = math.NaN()
		v.out[1] = math.NaN()
		return v.out
	}

	sumVMP := 0.0
	sumVMM := 0.0
	sumTR := 0.0
	for i := 0; i < v.period; i++ {
		sumVMP += v.vmPlus[i]
		sumVMM += v.vmMinus[i]
		sumTR += v.trValues[i]
	}

	if sumTR == 0 {
		v.out[0] = math.NaN()
		v.out[1] = math.NaN()
		return v.out
	}

	v.out[0] = sumVMP / sumTR
	v.out[1] = sumVMM / sumTR
	return v.out
}

func (v *Vortex) Reset() {
	v.vmPlus = make([]float64, v.period)
	v.vmMinus = make([]float64, v.period)
	v.trValues = make([]float64, v.period)
	v.prevHigh = 0
	v.prevLow = 0
	v.prevClose = 0
	v.index = 0
	v.count = 0
}

func (v *Vortex) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Vortex Indicator",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(v.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "VI+", Color: "#4CAF50", Style: indicators.StyleLine, Width: 2},
			{Name: "VI-", Color: "#F44336", Style: indicators.StyleLine, Width: 2},
		},
	}
}
