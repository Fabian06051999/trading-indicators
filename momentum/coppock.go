package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// CoppockCurve implements the Coppock Curve.
type CoppockCurve struct {
	wmaPeriod  int
	roc1Period int
	roc2Period int
	wma        *moving_averages.WMA
	buffer     []float64
	index      int
	count      int
	maxPeriod  int
}

func NewCoppockCurve(wmaPeriod, roc1Period, roc2Period int) *CoppockCurve {
	maxP := roc1Period
	if roc2Period > maxP {
		maxP = roc2Period
	}
	return &CoppockCurve{
		wmaPeriod:  wmaPeriod,
		roc1Period: roc1Period,
		roc2Period: roc2Period,
		wma:        moving_averages.NewWMA(wmaPeriod),
		buffer:     make([]float64, maxP+1),
		maxPeriod:  maxP,
	}
}

func (c *CoppockCurve) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		result[i] = c.UpdateAll(candle)[0]
	}
	return [][]float64{result}
}

func (c *CoppockCurve) UpdateAll(candle indicators.OHLCV) []float64 {
	c.buffer[c.index] = candle.Close
	c.count++

	if c.count <= c.maxPeriod {
		c.index = (c.index + 1) % (c.maxPeriod + 1)
		return []float64{math.NaN()}
	}

	// ROC calculations
	past1Idx := (c.index - c.roc1Period + c.maxPeriod + 1) % (c.maxPeriod + 1)
	past2Idx := (c.index - c.roc2Period + c.maxPeriod + 1) % (c.maxPeriod + 1)

	roc1 := 0.0
	roc2 := 0.0
	if c.buffer[past1Idx] != 0 {
		roc1 = ((candle.Close - c.buffer[past1Idx]) / c.buffer[past1Idx]) * 100
	}
	if c.buffer[past2Idx] != 0 {
		roc2 = ((candle.Close - c.buffer[past2Idx]) / c.buffer[past2Idx]) * 100
	}

	c.index = (c.index + 1) % (c.maxPeriod + 1)

	sum := roc1 + roc2
	return []float64{c.wma.UpdateAll(indicators.OHLCV{Close: sum})[0]}
}

func (c *CoppockCurve) Reset() {
	c.wma.Reset()
	c.buffer = make([]float64, c.maxPeriod+1)
	c.index = 0
	c.count = 0
}

func (c *CoppockCurve) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Coppock Curve",
		Parameters: []indicators.Parameter{
			{Name: "WMA Period", DefaultValue: float64(c.wmaPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "ROC1 Period", DefaultValue: float64(c.roc1Period), Min: 5, Max: 30, Step: 1},
			{Name: "ROC2 Period", DefaultValue: float64(c.roc2Period), Min: 5, Max: 30, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Coppock",
				Color: "#283593",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
