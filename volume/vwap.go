package volume

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// VWAP implements the Volume Weighted Average Price.
// Resets each session (call Reset() at session boundaries).
type VWAP struct {
	cumTPV float64
	cumVol float64
	count  int
}

func NewVWAP() *VWAP {
	return &VWAP{}
}

func (v *VWAP) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		result[i] = v.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (v *VWAP) UpdateAll(candle indicators.OHLCV) []float64 {
	v.count++
	tp := (candle.High + candle.Low + candle.Close) / 3.0
	v.cumTPV += tp * candle.Volume
	v.cumVol += candle.Volume

	if v.cumVol == 0 {
		return []float64{math.NaN()}
	}
	return []float64{v.cumTPV / v.cumVol}
}

func (v *VWAP) Reset() {
	v.cumTPV = 0
	v.cumVol = 0
	v.count = 0
}

func (v *VWAP) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Volume Weighted Average Price",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "VWAP", Color: "#9C27B0", Style: indicators.StyleLine, Width: 2},
		},
	}
}
