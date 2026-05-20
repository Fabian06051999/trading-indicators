package moving_averages

import (
	"github.com/Fabian06051999/trading-indicators"
	"math"
)

// ZLEMA implements the Zero Lag Exponential Moving Average.
type ZLEMA struct {
	period int
	lag    int
	ema    *EMA
	buffer []float64
	index  int
	count  int
	out    []float64
}

func NewZLEMA(period int) *ZLEMA {
	period = indicators.ClampMin(period, 2)
	lag := (period - 1) / 2
	return &ZLEMA{
		period: period,
		lag:    lag,
		ema:    NewEMA(period),
		buffer: make([]float64, lag+1),
		out:    make([]float64, 1),
	}
}

func (z *ZLEMA) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	z.Reset()

	for i, c := range candles {
		z.update(c)
		result[i] = z.out[0]
	}
	return [][]float64{result}
}

func (z *ZLEMA) UpdateAll(candle indicators.OHLCV) []float64 {
	z.update(candle)
	return z.out
}

func (z *ZLEMA) update(candle indicators.OHLCV) {
	z.count++

	// Store current close in buffer
	z.buffer[z.index] = candle.Close
	z.index = (z.index + 1) % (z.lag + 1)

	if z.count <= z.lag {
		z.out[0] = math.NaN()
		return
	}

	// Get lagged value
	laggedValue := z.buffer[z.index]
	// Zero-lag adjusted price
	adjustedPrice := 2*candle.Close - laggedValue

	z.out[0] = z.ema.UpdateAll(indicators.OHLCV{Close: adjustedPrice})[0]
}

func (z *ZLEMA) Reset() {
	z.ema.Reset()
	z.buffer = make([]float64, z.lag+1)
	z.index = 0
	z.count = 0
}

func (z *ZLEMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Zero Lag EMA",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(z.period), Min: 1, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "ZLEMA", Color: "#673AB7", Style: indicators.StyleLine, Width: 2},
		},
	}
}
