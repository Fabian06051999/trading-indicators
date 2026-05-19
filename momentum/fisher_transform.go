package momentum

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// FisherTransform implements the Ehlers Fisher Transform.
type FisherTransform struct {
	period    int
	highs     []float64
	lows      []float64
	index     int
	count     int
	value     float64
	prevValue float64
	prevNorm  float64
}

func NewFisherTransform(period int) *FisherTransform {
	return &FisherTransform{
		period: period,
		highs:  make([]float64, period),
		lows:   make([]float64, period),
	}
}

func (f *FisherTransform) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	fisher := make([]float64, len(candles))
	trigger := make([]float64, len(candles))
	f.Reset()

	for i, c := range candles {
		values := f.UpdateAll(c)
		fisher[i] = values[0]
		trigger[i] = values[1]
	}
	return [][]float64{fisher, trigger}
}

func (f *FisherTransform) UpdateAll(candle indicators.OHLCV) []float64 {
	f.highs[f.index] = candle.High
	f.lows[f.index] = candle.Low
	f.index = (f.index + 1) % f.period
	if f.count < f.period {
		f.count++
	}
	if f.count < f.period {
		return []float64{math.NaN(), math.NaN()}
	}

	// Find highest high and lowest low
	hh := f.highs[0]
	ll := f.lows[0]
	for i := 1; i < f.period; i++ {
		if f.highs[i] > hh {
			hh = f.highs[i]
		}
		if f.lows[i] < ll {
			ll = f.lows[i]
		}
	}

	// Normalize price to -1..1
	midPrice := (candle.High + candle.Low) / 2.0
	norm := 0.0
	if hh-ll != 0 {
		norm = 2*((midPrice-ll)/(hh-ll)) - 1
	}

	// Smooth
	norm = 0.33*norm + 0.67*f.prevNorm
	if norm > 0.99 {
		norm = 0.99
	}
	if norm < -0.99 {
		norm = -0.99
	}
	f.prevNorm = norm

	// Fisher Transform
	f.prevValue = f.value
	f.value = 0.5*math.Log((1+norm)/(1-norm)) + 0.5*f.value

	return []float64{f.value, f.prevValue}
}

func (f *FisherTransform) Reset() {
	f.highs = make([]float64, f.period)
	f.lows = make([]float64, f.period)
	f.index = 0
	f.count = 0
	f.value = 0
	f.prevValue = 0
	f.prevNorm = 0
}

func (f *FisherTransform) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Fisher Transform",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(f.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Fisher",
				Color: "#1B5E20",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
			{Name: "Trigger", Color: "#F44336", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
