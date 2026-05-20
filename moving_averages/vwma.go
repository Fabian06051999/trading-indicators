package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// VWMA implements the Volume Weighted Moving Average.
type VWMA struct {
	period   int
	priceBuf []float64
	volBuf   []float64
	index    int
	count    int
	priceVol float64
	volSum   float64
	out      []float64
}

func NewVWMA(period int) *VWMA {
	period = indicators.ClampMin(period, 1)
	return &VWMA{
		period:   period,
		priceBuf: make([]float64, period),
		volBuf:   make([]float64, period),
		out:      make([]float64, 1),
	}
}

func (v *VWMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		v.update(c)
		result[i] = v.out[0]
	}
	return [][]float64{result}
}

func (v *VWMA) UpdateAll(candle indicators.OHLCV) []float64 {
	v.update(candle)
	return v.out
}

func (v *VWMA) update(candle indicators.OHLCV) {
	oldPV := v.priceBuf[v.index] * v.volBuf[v.index]
	oldVol := v.volBuf[v.index]

	v.priceBuf[v.index] = candle.Close
	v.volBuf[v.index] = candle.Volume

	v.priceVol += candle.Close*candle.Volume - oldPV
	v.volSum += candle.Volume - oldVol

	v.index = (v.index + 1) % v.period
	if v.count < v.period {
		v.count++
	}
	if v.count < v.period {
		v.out[0] = math.NaN()
		return
	}
	if v.volSum == 0 {
		v.out[0] = math.NaN()
		return
	}
	v.out[0] = v.priceVol / v.volSum
}

func (v *VWMA) Reset() {
	v.priceBuf = make([]float64, v.period)
	v.volBuf = make([]float64, v.period)
	v.index = 0
	v.count = 0
	v.priceVol = 0
	v.volSum = 0
}

func (v *VWMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Volume Weighted Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(v.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "VWMA", Color: "#607D8B", Style: indicators.StyleLine, Width: 2},
		},
	}
}
