package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// KAMA implements the Kaufman Adaptive Moving Average.
type KAMA struct {
	period     int
	fastPeriod int
	slowPeriod int
	buffer     []float64
	value      float64
	index      int
	count      int
}

func NewKAMA(period, fastPeriod, slowPeriod int) *KAMA {
	return &KAMA{
		period:     period,
		fastPeriod: fastPeriod,
		slowPeriod: slowPeriod,
		buffer:     make([]float64, period+1),
	}
}

func (k *KAMA) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	k.Reset()

	for i, c := range candles {
		result[i] = k.Update(c)
	}
	return result
}

func (k *KAMA) Update(candle indicators.OHLCV) float64 {
	k.buffer[k.index] = candle.Close
	k.count++

	if k.count <= k.period {
		k.index = (k.index + 1) % (k.period + 1)
		if k.count == k.period {
			k.value = candle.Close
			return k.value
		}
		return math.NaN()
	}

	oldestIdx := (k.index + 1) % (k.period + 1)

	// Direction: |close - close[period]|
	direction := math.Abs(candle.Close - k.buffer[oldestIdx])

	// Volatility: sum of |close[i] - close[i-1]|
	volatility := 0.0
	for i := 1; i <= k.period; i++ {
		curr := (oldestIdx + i) % (k.period + 1)
		prev := (oldestIdx + i - 1) % (k.period + 1)
		volatility += math.Abs(k.buffer[curr] - k.buffer[prev])
	}

	// Efficiency Ratio
	er := 0.0
	if volatility != 0 {
		er = direction / volatility
	}

	// Smoothing constant
	fastSC := 2.0 / float64(k.fastPeriod+1)
	slowSC := 2.0 / float64(k.slowPeriod+1)
	sc := er*(fastSC-slowSC) + slowSC
	sc = sc * sc

	k.value = k.value + sc*(candle.Close-k.value)
	k.index = (k.index + 1) % (k.period + 1)
	return k.value
}

func (k *KAMA) Reset() {
	k.buffer = make([]float64, k.period+1)
	k.value = 0
	k.index = 0
	k.count = 0
}

func (k *KAMA) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Kaufman Adaptive Moving Average",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(k.period), Min: 2, Max: 500, Step: 1},
			{Name: "Fast Period", DefaultValue: float64(k.fastPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "Slow Period", DefaultValue: float64(k.slowPeriod), Min: 10, Max: 100, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "KAMA", Color: "#FF5722", Style: indicators.StyleLine, Width: 2},
		},
	}
}
