package volatility

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// DonchianChannel implements Donchian Channels (highest high / lowest low).
type DonchianChannel struct {
	out    []float64
	period int
	highs  []float64
	lows   []float64
	index  int
	count  int
}

func NewDonchianChannel(period int) *DonchianChannel {
	return &DonchianChannel{
		period: period,
		highs:  make([]float64, period),
		lows:   make([]float64, period),
		out:    make([]float64, 3),
	}
}

func (d *DonchianChannel) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	upper := make([]float64, len(candles))
	middle := make([]float64, len(candles))
	lower := make([]float64, len(candles))
	d.Reset()

	for i, c := range candles {
		values := d.UpdateAll(c)
		upper[i] = values[0]
		middle[i] = values[1]
		lower[i] = values[2]
	}
	return [][]float64{upper, middle, lower}
}

func (d *DonchianChannel) UpdateAll(candle indicators.OHLCV) []float64 {
	d.highs[d.index] = candle.High
	d.lows[d.index] = candle.Low
	d.index = (d.index + 1) % d.period
	if d.count < d.period {
		d.count++
	}
	if d.count < d.period {
		d.out[0] = math.NaN()
		d.out[1] = math.NaN()
		d.out[2] = math.NaN()
		return d.out
	}

	hh := d.highs[0]
	ll := d.lows[0]
	for i := 1; i < d.period; i++ {
		if d.highs[i] > hh {
			hh = d.highs[i]
		}
		if d.lows[i] < ll {
			ll = d.lows[i]
		}
	}

	mid := (hh + ll) / 2.0
	d.out[0] = hh
	d.out[1] = mid
	d.out[2] = ll
	return d.out
}

func (d *DonchianChannel) Reset() {
	d.highs = make([]float64, d.period)
	d.lows = make([]float64, d.period)
	d.index = 0
	d.count = 0
}

func (d *DonchianChannel) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Donchian Channel",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(d.period), Min: 2, Max: 200, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Upper", Color: "#EF5350", FillColor: "#EF535015", Style: indicators.StyleLine, Width: 1, Opacity: 1.0},
			{Name: "Middle", Color: "#9E9E9E", Style: indicators.StyleDashed, Width: 1, Opacity: 0.7, DashArray: []int{5, 5}},
			{Name: "Lower", Color: "#66BB6A", FillColor: "#66BB6A15", Style: indicators.StyleLine, Width: 1, Opacity: 1.0},
		},
	}
}
