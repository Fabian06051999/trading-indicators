package trend

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// EMAEnvelope implements Moving Average Envelopes.
type EMAEnvelope struct {
	period     int
	percentage float64
	ema        *moving_averages.EMA
}

func NewEMAEnvelope(period int, percentage float64) *EMAEnvelope {
	return &EMAEnvelope{
		period:     period,
		percentage: percentage,
		ema:        moving_averages.NewEMA(period),
	}
}

func (e *EMAEnvelope) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	upper := make([]float64, len(candles))
	middle := make([]float64, len(candles))
	lower := make([]float64, len(candles))
	e.Reset()

	for i, c := range candles {
		values := e.UpdateAll(c)
		upper[i] = values[0]
		middle[i] = values[1]
		lower[i] = values[2]
	}
	return [][]float64{upper, middle, lower}
}

func (e *EMAEnvelope) UpdateAll(candle indicators.OHLCV) []float64 {
	mid := e.ema.Update(candle)
	if mid == 0 {
		return []float64{math.NaN(), math.NaN(), math.NaN()}
	}

	offset := mid * (e.percentage / 100.0)
	return []float64{mid + offset, mid, mid - offset}
}

func (e *EMAEnvelope) Reset() {
	e.ema.Reset()
}

func (e *EMAEnvelope) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Moving Average Envelope",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(e.period), Min: 2, Max: 200, Step: 1},
			{Name: "Percentage", DefaultValue: e.percentage, Min: 0.1, Max: 10, Step: 0.1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Upper", Color: "#EF5350", Style: indicators.StyleLine, Width: 1},
			{Name: "Middle", Color: "#2196F3", Style: indicators.StyleDashed, Width: 1},
			{Name: "Lower", Color: "#66BB6A", Style: indicators.StyleLine, Width: 1},
		},
	}
}
