package trend

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// DPO implements the Detrended Price Oscillator.
type DPO struct {
	period int
	shift  int
	buffer []float64
	sum    float64
	index  int
	count  int
	out    []float64
}

func NewDPO(period int) *DPO {
	shift := period/2 + 1
	return &DPO{
		period: period,
		shift:  shift,
		buffer: make([]float64, period+shift),
		out:    make([]float64, 1),
	}
}

func (d *DPO) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		d.update(c)
		result[i] = d.out[0]
	}
	return [][]float64{result}
}

func (d *DPO) UpdateAll(candle indicators.OHLCV) []float64 {
	d.update(candle)
	return d.out
}

func (d *DPO) update(candle indicators.OHLCV) {
	bufLen := len(d.buffer)
	d.buffer[d.index] = candle.Close
	d.count++

	if d.count < d.period+d.shift {
		d.index = (d.index + 1) % bufLen
		d.out[0] = math.NaN()
		return
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
	d.out[0] = dpoVal
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
