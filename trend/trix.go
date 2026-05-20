package trend

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
	"math"
)

// TRIX implements the Triple Exponential Average (1-day ROC of triple EMA).
type TRIX struct {
	period int
	ema1   *moving_averages.EMA
	ema2   *moving_averages.EMA
	ema3   *moving_averages.EMA
	prev   float64
	count  int
	out    []float64
}

func NewTRIX(period int) *TRIX {
	return &TRIX{
		period: period,
		ema1:   moving_averages.NewEMA(period),
		ema2:   moving_averages.NewEMA(period),
		ema3:   moving_averages.NewEMA(period),
		out:    make([]float64, 1),
	}
}

func (t *TRIX) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		t.update(c)
		result[i] = t.out[0]
	}
	return [][]float64{result}
}

func (t *TRIX) UpdateAll(candle indicators.OHLCV) []float64 {
	t.update(candle)
	return t.out
}

func (t *TRIX) update(candle indicators.OHLCV) {
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
		t.prev = e3
		t.out[0] = math.NaN()
		return
	}

	if math.IsNaN(t.prev) {
		t.prev = e3
		t.out[0] = math.NaN()
		return
	}

	t.out[0] = ((e3 - t.prev) / t.prev) * 100
	t.prev = e3
}

func (t *TRIX) Reset() {
	t.ema1.Reset()
	t.ema2.Reset()
	t.ema3.Reset()
	t.prev = 0
	t.count = 0
}

func (t *TRIX) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "TRIX",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(t.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "TRIX",
				Color: "#7B1FA2",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
