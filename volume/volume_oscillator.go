package volume

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
	"math"
)

// VolumeOscillator implements the Volume Oscillator (fast EMA - slow EMA of volume).
type VolumeOscillator struct {
	out        []float64
	fastPeriod int
	slowPeriod int
	fastEMA    *moving_averages.EMA
	slowEMA    *moving_averages.EMA
}

func NewVolumeOscillator(fastPeriod, slowPeriod int) *VolumeOscillator {
	return &VolumeOscillator{
		fastPeriod: fastPeriod,
		slowPeriod: slowPeriod,
		fastEMA:    moving_averages.NewEMA(fastPeriod),
		slowEMA:    moving_averages.NewEMA(slowPeriod),
		out:        make([]float64, 1),
	}
}

func (v *VolumeOscillator) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		result[i] = v.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (v *VolumeOscillator) UpdateAll(candle indicators.OHLCV) []float64 {
	volCandle := indicators.OHLCV{Close: candle.Volume}
	fast := v.fastEMA.UpdateAll(volCandle)[0]
	slow := v.slowEMA.UpdateAll(volCandle)[0]

	if fast == 0 || slow == 0 {
		v.out[0] = math.NaN()
		return v.out
	}

	v.out[0] = ((fast - slow) / slow) * 100
	return v.out
}

func (v *VolumeOscillator) Reset() {
	v.fastEMA.Reset()
	v.slowEMA.Reset()
}

func (v *VolumeOscillator) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Volume Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Fast Period", DefaultValue: float64(v.fastPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "Slow Period", DefaultValue: float64(v.slowPeriod), Min: 5, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "VO",
				Color: "#558B2F",
				Style: indicators.StyleHistogram,
				Width: 1,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
