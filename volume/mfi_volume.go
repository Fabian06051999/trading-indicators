package volume

import (
	"github.com/Fabian06051999/trading-indicators"
)

// VolumeProfile provides volume-at-price distribution for a given window.
type VolumeProfile struct {
	period int
	bins   int
	buffer []indicators.OHLCV
	index  int
	count  int
}

func NewVolumeProfile(period, bins int) *VolumeProfile {
	return &VolumeProfile{
		period: period,
		bins:   bins,
		buffer: make([]indicators.OHLCV, period),
	}
}

func (vp *VolumeProfile) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	vp.Reset()

	for i, c := range candles {
		result[i] = vp.Update(c)
	}
	return result
}

// Update returns the POC (Point of Control) price level.
func (vp *VolumeProfile) Update(candle indicators.OHLCV) float64 {
	vp.buffer[vp.index] = candle
	vp.index = (vp.index + 1) % vp.period
	if vp.count < vp.period {
		vp.count++
	}
	if vp.count < vp.period {
		return 0
	}

	// Find price range
	high := vp.buffer[0].High
	low := vp.buffer[0].Low
	for i := 1; i < vp.period; i++ {
		if vp.buffer[i].High > high {
			high = vp.buffer[i].High
		}
		if vp.buffer[i].Low < low {
			low = vp.buffer[i].Low
		}
	}

	if high == low {
		return high
	}

	binSize := (high - low) / float64(vp.bins)
	volumes := make([]float64, vp.bins)

	for i := 0; i < vp.period; i++ {
		c := vp.buffer[i]
		tp := (c.High + c.Low + c.Close) / 3.0
		bin := int((tp - low) / binSize)
		if bin >= vp.bins {
			bin = vp.bins - 1
		}
		if bin < 0 {
			bin = 0
		}
		volumes[bin] += c.Volume
	}

	// Find POC (bin with highest volume)
	maxBin := 0
	maxVol := volumes[0]
	for i := 1; i < vp.bins; i++ {
		if volumes[i] > maxVol {
			maxVol = volumes[i]
			maxBin = i
		}
	}

	return low + (float64(maxBin)+0.5)*binSize
}

func (vp *VolumeProfile) Reset() {
	vp.buffer = make([]indicators.OHLCV, vp.period)
	vp.index = 0
	vp.count = 0
}

func (vp *VolumeProfile) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Volume Profile (POC)",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(vp.period), Min: 10, Max: 500, Step: 10},
			{Name: "Bins", DefaultValue: float64(vp.bins), Min: 10, Max: 100, Step: 5},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "POC", Color: "#FF6F00", Style: indicators.StyleDashed, Width: 2},
		},
	}
}
