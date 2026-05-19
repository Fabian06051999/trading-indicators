package volume

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// ChaikinMF implements the Chaikin Money Flow indicator.
type ChaikinMF struct {
	period   int
	mfvBuf   []float64
	volBuf   []float64
	index    int
	count    int
}

func NewChaikinMF(period int) *ChaikinMF {
	return &ChaikinMF{
		period: period,
		mfvBuf: make([]float64, period),
		volBuf: make([]float64, period),
	}
}

func (c *ChaikinMF) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		result[i] = c.Update(candle)
	}
	return result
}

func (c *ChaikinMF) Update(candle indicators.OHLCV) float64 {
	hl := candle.High - candle.Low
	mfm := 0.0
	if hl != 0 {
		mfm = ((candle.Close - candle.Low) - (candle.High - candle.Close)) / hl
	}
	mfv := mfm * candle.Volume

	c.mfvBuf[c.index] = mfv
	c.volBuf[c.index] = candle.Volume
	c.index = (c.index + 1) % c.period
	if c.count < c.period {
		c.count++
	}
	if c.count < c.period {
		return math.NaN()
	}

	sumMFV := 0.0
	sumVol := 0.0
	for i := 0; i < c.period; i++ {
		sumMFV += c.mfvBuf[i]
		sumVol += c.volBuf[i]
	}

	if sumVol == 0 {
		return math.NaN()
	}
	return sumMFV / sumVol
}

func (c *ChaikinMF) Reset() {
	c.mfvBuf = make([]float64, c.period)
	c.volBuf = make([]float64, c.period)
	c.index = 0
	c.count = 0
}

func (c *ChaikinMF) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Chaikin Money Flow",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(c.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "CMF",
				Color:  "#00838F",
				Style:  indicators.StyleHistogram,
				Width:  1,
				YRange: &indicators.YRange{Min: -1, Max: 1},
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
