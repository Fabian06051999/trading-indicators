package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// T3 implements the Tillson T3 Moving Average.
type T3 struct {
	period int
	vf     float64
	ema1   *EMA
	ema2   *EMA
	ema3   *EMA
	ema4   *EMA
	ema5   *EMA
	ema6   *EMA
	c1     float64
	c2     float64
	c3     float64
	c4     float64
	count  int
	out    []float64
}

func NewT3(period int, volumeFactor float64) *T3 {
	vf := volumeFactor
	c1 := -(vf * vf * vf)
	c2 := 3*vf*vf + 3*vf*vf*vf
	c3 := -6*vf*vf - 3*vf - 3*vf*vf*vf
	c4 := 1 + 3*vf + vf*vf*vf + 3*vf*vf

	return &T3{
		period: period,
		vf:     vf,
		ema1:   NewEMA(period),
		ema2:   NewEMA(period),
		ema3:   NewEMA(period),
		ema4:   NewEMA(period),
		ema5:   NewEMA(period),
		ema6:   NewEMA(period),
		c1:     c1,
		c2:     c2,
		c3:     c3,
		c4:     c4,
		out:    make([]float64, 1),
	}
}

func (t *T3) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		t.update(c)
		result[i] = t.out[0]
	}
	return [][]float64{result}
}

func (t *T3) UpdateAll(candle indicators.OHLCV) []float64 {
	t.update(candle)
	return t.out
}

func (t *T3) update(candle indicators.OHLCV) {
	t.count++
	e1 := t.ema1.UpdateAll(candle)[0]
	if math.IsNaN(e1) {
		t.out[0] = math.NaN()
		return
	}
	e2 := t.ema2.UpdateAll(indicators.OHLCV{Close: e1})[0]
	if math.IsNaN(e2) {
		t.out[0] = math.NaN()
		return
	}
	e3 := t.ema3.UpdateAll(indicators.OHLCV{Close: e2})[0]
	if math.IsNaN(e3) {
		t.out[0] = math.NaN()
		return
	}
	e4 := t.ema4.UpdateAll(indicators.OHLCV{Close: e3})[0]
	if math.IsNaN(e4) {
		t.out[0] = math.NaN()
		return
	}
	e5 := t.ema5.UpdateAll(indicators.OHLCV{Close: e4})[0]
	if math.IsNaN(e5) {
		t.out[0] = math.NaN()
		return
	}
	e6 := t.ema6.UpdateAll(indicators.OHLCV{Close: e5})[0]
	if math.IsNaN(e6) {
		t.out[0] = math.NaN()
		return
	}

	t.out[0] = t.c1*e6 + t.c2*e5 + t.c3*e4 + t.c4*e3
}

func (t *T3) Reset() {
	t.ema1.Reset()
	t.ema2.Reset()
	t.ema3.Reset()
	t.ema4.Reset()
	t.ema5.Reset()
	t.ema6.Reset()
	t.count = 0
}

func (t *T3) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Tillson T3",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(t.period), Min: 2, Max: 200, Step: 1},
			{Name: "Volume Factor", DefaultValue: t.vf, Min: 0, Max: 1, Step: 0.1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "T3", Color: "#FF6F00", Style: indicators.StyleLine, Width: 2},
		},
	}
}
