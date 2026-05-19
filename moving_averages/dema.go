package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
)

// DEMA implements the Double Exponential Moving Average.
type DEMA struct {
	period int
	ema1   *EMA
	ema2   *EMA
	count  int
}

func NewDEMA(period int) *DEMA {
	return &DEMA{
		period: period,
		ema1:   NewEMA(period),
		ema2:   NewEMA(period),
	}
}

func (d *DEMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		result[i] = d.Update(c)
	}
	return result
}

func (d *DEMA) Update(candle indicators.OHLCV) float64 {
	d.count++
	ema1Val := d.ema1.Update(candle)
	if ema1Val == 0 {
		return 0
	}

	ema2Val := d.ema2.Update(indicators.OHLCV{Close: ema1Val})
	if ema2Val == 0 {
		return 0
	}

	return 2*ema1Val - ema2Val
}

func (d *DEMA) Reset() {
	d.ema1.Reset()
	d.ema2.Reset()
	d.count = 0
}

func (d *DEMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Double Exponential Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(d.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "DEMA", Color: "#9C27B0", Style: indicators.StyleLine, Width: 2},
		},
	}
}
