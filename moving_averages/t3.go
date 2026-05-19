package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
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
	}
}

func (t *T3) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		result[i] = t.Update(c)
	}
	return result
}

func (t *T3) Update(candle indicators.OHLCV) float64 {
	t.count++
	e1 := t.ema1.Update(candle)
	if e1 == 0 {
		return 0
	}
	e2 := t.ema2.Update(indicators.OHLCV{Close: e1})
	if e2 == 0 {
		return 0
	}
	e3 := t.ema3.Update(indicators.OHLCV{Close: e2})
	if e3 == 0 {
		return 0
	}
	e4 := t.ema4.Update(indicators.OHLCV{Close: e3})
	if e4 == 0 {
		return 0
	}
	e5 := t.ema5.Update(indicators.OHLCV{Close: e4})
	if e5 == 0 {
		return 0
	}
	e6 := t.ema6.Update(indicators.OHLCV{Close: e5})
	if e6 == 0 {
		return 0
	}

	return t.c1*e6 + t.c2*e5 + t.c3*e4 + t.c4*e3
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
