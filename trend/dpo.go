package trend

import (
	"github.com/Fabian06051999/trading-indicators"
)

// DPO implements the Detrended Price Oscillator.
type DPO struct {
	period int
	shift  int
	buffer []float64
	sum    float64
	index  int
	count  int
}

func NewDPO(period int) *DPO {
	shift := period/2 + 1
	return &DPO{
		period: period,
		shift:  shift,
		buffer: make([]float64, period+shift),
	}
}

func (d *DPO) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		result[i] = d.Update(c)
	}
	return result
}

func (d *DPO) Update(candle indicators.OHLCV) float64 {
	bufLen := len(d.buffer)
	d.buffer[d.index] = candle.Close
	d.count++

	if d.count < d.period+d.shift {
		d.index = (d.index + 1) % bufLen
		return 0
	}

	// SMA of current period window
	smaStart := (d.index - d.period - d.shift + 1 + bufLen) % bufLen
	sum := 0.0
	for i := 0; i < d.period; i++ {
		idx := (smaStart + i) % bufLen
		sum += d.buffer[idx]
	}
	sma := sum / float64(d.period)

	// Price shifted back
	shiftedIdx := (d.index - d.shift + 1 + bufLen) % bufLen
	dpoVal := d.buffer[shiftedIdx] - sma

	d.index = (d.index + 1) % bufLen
	return dpoVal
}

func (d *DPO) Reset() {
	d.buffer = make([]float64, d.period+d.shift)
	d.sum = 0
	d.index = 0
	d.count = 0
}

func (d *DPO) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Detrended Price Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(d.period), Min: 2, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "DPO",
				Color: "#00695C",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
