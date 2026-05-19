package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// VIDYA implements the Variable Index Dynamic Average (Chande).
type VIDYA struct {
	period    int
	cmoPeriod int
	sc        float64
	value     float64
	gains     []float64
	losses    []float64
	prevClose float64
	index     int
	count     int
}

func NewVIDYA(period, cmoPeriod int) *VIDYA {
	return &VIDYA{
		period:    period,
		cmoPeriod: cmoPeriod,
		sc:        2.0 / float64(period+1),
		gains:     make([]float64, cmoPeriod),
		losses:    make([]float64, cmoPeriod),
	}
}

func (v *VIDYA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	v.Reset()

	for i, c := range candles {
		result[i] = v.Update(c)
	}
	return result
}

func (v *VIDYA) Update(candle indicators.OHLCV) float64 {
	v.count++

	if v.count == 1 {
		v.prevClose = candle.Close
		v.value = candle.Close
		return 0
	}

	change := candle.Close - v.prevClose
	v.prevClose = candle.Close

	gain := 0.0
	loss := 0.0
	if change > 0 {
		gain = change
	} else {
		loss = -change
	}

	v.gains[v.index] = gain
	v.losses[v.index] = loss
	v.index = (v.index + 1) % v.cmoPeriod

	if v.count <= v.cmoPeriod {
		return 0
	}

	// CMO = (sumGains - sumLosses) / (sumGains + sumLosses)
	sumGains := 0.0
	sumLosses := 0.0
	for i := 0; i < v.cmoPeriod; i++ {
		sumGains += v.gains[i]
		sumLosses += v.losses[i]
	}

	cmo := 0.0
	if sumGains+sumLosses != 0 {
		cmo = math.Abs((sumGains - sumLosses) / (sumGains + sumLosses))
	}

	v.value = v.value + v.sc*cmo*(candle.Close-v.value)
	return v.value
}

func (v *VIDYA) Reset() {
	v.value = 0
	v.gains = make([]float64, v.cmoPeriod)
	v.losses = make([]float64, v.cmoPeriod)
	v.prevClose = 0
	v.index = 0
	v.count = 0
}

func (v *VIDYA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Variable Index Dynamic Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(v.period), Min: 2, Max: 500, Step: 1},
			{Name: "CMO Period", DefaultValue: float64(v.cmoPeriod), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "VIDYA", Color: "#5C6BC0", Style: indicators.StyleLine, Width: 2},
		},
	}
}
