package volume

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// ForceIndex implements the Force Index (Elder).
type ForceIndex struct {
	period    int
	ema       *moving_averages.EMA
	prevClose float64
	count     int
}

func NewForceIndex(period int) *ForceIndex {
	return &ForceIndex{
		period: period,
		ema:    moving_averages.NewEMA(period),
	}
}

func (f *ForceIndex) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	f.Reset()

	for i, c := range candles {
		result[i] = f.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (f *ForceIndex) UpdateAll(candle indicators.OHLCV) []float64 {
	f.count++

	if f.count == 1 {
		f.prevClose = candle.Close
		return []float64{math.NaN()}
	}

	rawForce := (candle.Close - f.prevClose) * candle.Volume
	f.prevClose = candle.Close

	return []float64{f.ema.UpdateAll(indicators.OHLCV{Close: rawForce})[0]}
}

func (f *ForceIndex) Reset() {
	f.ema.Reset()
	f.prevClose = 0
	f.count = 0
}

func (f *ForceIndex) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Force Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(f.period), Min: 1, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "FI",
				Color: "#1B5E20",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
