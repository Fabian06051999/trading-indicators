package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// FRAMA implements the Fractal Adaptive Moving Average.
type FRAMA struct {
	period int
	value  float64
	buffer []float64
	index  int
	count  int
}

func NewFRAMA(period int) *FRAMA {
	period = indicators.ClampMin(period, 4)
	return &FRAMA{
		period: period,
		buffer: make([]float64, period),
	}
}

func (f *FRAMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	f.Reset()

	for i, c := range candles {
		result[i] = f.Update(c)
	}
	return result
}

func (f *FRAMA) Update(candle indicators.OHLCV) float64 {
	f.buffer[f.index] = candle.Close
	f.index = (f.index + 1) % f.period
	f.count++

	if f.count < f.period {
		return math.NaN()
	}

	if f.count == f.period {
		f.value = candle.Close
		return f.value
	}

	half := f.period / 2

	// Find highest/lowest for first half, second half, and full period
	hh1, ll1 := f.rangeHL(0, half)
	hh2, ll2 := f.rangeHL(half, f.period)
	hh3, ll3 := f.rangeHL(0, f.period)

	n1 := (hh1 - ll1) / float64(half)
	n2 := (hh2 - ll2) / float64(half)
	n3 := (hh3 - ll3) / float64(f.period)

	dim := 0.0
	if n1+n2 > 0 && n3 > 0 {
		dim = (math.Log(n1+n2) - math.Log(n3)) / math.Log(2)
	}

	alpha := math.Exp(-4.6 * (dim - 1))
	if alpha > 1 {
		alpha = 1
	}
	if alpha < 0.01 {
		alpha = 0.01
	}

	f.value = alpha*candle.Close + (1-alpha)*f.value
	return f.value
}

func (f *FRAMA) rangeHL(start, end int) (float64, float64) {
	hh := -math.MaxFloat64
	ll := math.MaxFloat64
	for i := start; i < end; i++ {
		idx := (f.index + i) % f.period
		if f.buffer[idx] > hh {
			hh = f.buffer[idx]
		}
		if f.buffer[idx] < ll {
			ll = f.buffer[idx]
		}
	}
	return hh, ll
}

func (f *FRAMA) Reset() {
	f.buffer = make([]float64, f.period)
	f.value = 0
	f.index = 0
	f.count = 0
}

func (f *FRAMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Fractal Adaptive Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(f.period), Min: 4, Max: 500, Step: 2},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "FRAMA", Color: "#AB47BC", Style: indicators.StyleLine, Width: 2},
		},
	}
}
