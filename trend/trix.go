package trend

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// TRIX implements the Triple Exponential Average (1-day ROC of triple EMA).
type TRIX struct {
	period int
	ema1   *moving_averages.EMA
	ema2   *moving_averages.EMA
	ema3   *moving_averages.EMA
	prev   float64
	count  int
}

func NewTRIX(period int) *TRIX {
	return &TRIX{
		period: period,
		ema1:   moving_averages.NewEMA(period),
		ema2:   moving_averages.NewEMA(period),
		ema3:   moving_averages.NewEMA(period),
	}
}

func (t *TRIX) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		result[i] = t.Update(c)
	}
	return result
}

func (t *TRIX) Update(candle indicators.OHLCV) float64 {
	t.count++
	e1 := t.ema1.Update(candle)
	if e1 == 0 {
		return math.NaN()
	}
	e2 := t.ema2.Update(indicators.OHLCV{Close: e1})
	if e2 == 0 {
		return math.NaN()
	}
	e3 := t.ema3.Update(indicators.OHLCV{Close: e2})
	if e3 == 0 {
		t.prev = e3
		return math.NaN()
	}

	if t.prev == 0 {
		t.prev = e3
		return math.NaN()
	}

	result := ((e3 - t.prev) / t.prev) * 100
	t.prev = e3
	return result
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
