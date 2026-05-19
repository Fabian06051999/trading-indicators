package moving_averages

import (
	"math"
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
	period = indicators.ClampMin(period, 1)
	return &DEMA{
		period: period,
		ema1:   NewEMA(period),
		ema2:   NewEMA(period),
	}
}

func (d *DEMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		result[i] = d.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (d *DEMA) UpdateAll(candle indicators.OHLCV) []float64 {
	d.count++
	ema1Val := d.ema1.UpdateAll(candle)[0]
	if ema1Val == 0 {
		return []float64{math.NaN()}
	}

	ema2Val := d.ema2.UpdateAll(indicators.OHLCV{Close: ema1Val})[0]
	if ema2Val == 0 {
		return []float64{math.NaN()}
	}

	return []float64{2*ema1Val - ema2Val}
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
