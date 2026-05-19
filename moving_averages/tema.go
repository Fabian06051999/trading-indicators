package moving_averages

import (
	"math"
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
	period = indicators.ClampMin(period, 1)
	return &TEMA{
		period: period,
		ema1:   NewEMA(period),
		ema2:   NewEMA(period),
		ema3:   NewEMA(period),
	}
}

func (t *TEMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		result[i] = t.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (t *TEMA) UpdateAll(candle indicators.OHLCV) []float64 {
	t.count++
	ema1Val := t.ema1.UpdateAll(candle)[0]
	if ema1Val == 0 {
		return []float64{math.NaN()}
	}

	ema2Val := t.ema2.UpdateAll(indicators.OHLCV{Close: ema1Val})[0]
	if ema2Val == 0 {
		return []float64{math.NaN()}
	}

	ema3Val := t.ema3.UpdateAll(indicators.OHLCV{Close: ema2Val})[0]
	if ema3Val == 0 {
		return []float64{math.NaN()}
	}

	return []float64{3*ema1Val - 3*ema2Val + ema3Val}
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
