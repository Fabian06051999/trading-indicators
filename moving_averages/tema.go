package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
)

// TEMA implements the Triple Exponential Moving Average.
type TEMA struct {
	period int
	ema1   *EMA
	ema2   *EMA
	ema3   *EMA
	count  int
}

func NewTEMA(period int) *TEMA {
	return &TEMA{
		period: period,
		ema1:   NewEMA(period),
		ema2:   NewEMA(period),
		ema3:   NewEMA(period),
	}
}

func (t *TEMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		result[i] = t.Update(c)
	}
	return result
}

func (t *TEMA) Update(candle indicators.OHLCV) float64 {
	t.count++
	ema1Val := t.ema1.Update(candle)
	if ema1Val == 0 {
		return 0
	}

	ema2Val := t.ema2.Update(indicators.OHLCV{Close: ema1Val})
	if ema2Val == 0 {
		return 0
	}

	ema3Val := t.ema3.Update(indicators.OHLCV{Close: ema2Val})
	if ema3Val == 0 {
		return 0
	}

	return 3*ema1Val - 3*ema2Val + ema3Val
}

func (t *TEMA) Reset() {
	t.ema1.Reset()
	t.ema2.Reset()
	t.ema3.Reset()
	t.count = 0
}

func (t *TEMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Triple Exponential Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(t.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "TEMA", Color: "#E91E63", Style: indicators.StyleLine, Width: 2},
		},
	}
}
