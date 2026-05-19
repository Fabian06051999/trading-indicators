package volume

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// VolumeOscillator implements the Volume Oscillator (fast EMA - slow EMA of volume).
type VolumeOscillator struct {
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
	}
}

func (v *VolumeOscillator) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		result[i] = v.Update(c)
	}
	return result
}

func (v *VolumeOscillator) Update(candle indicators.OHLCV) float64 {
	volCandle := indicators.OHLCV{Close: candle.Volume}
	fast := v.fastEMA.Update(volCandle)
	slow := v.slowEMA.Update(volCandle)

	if fast == 0 || slow == 0 {
		return 0
	}

	return ((fast - slow) / slow) * 100
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
