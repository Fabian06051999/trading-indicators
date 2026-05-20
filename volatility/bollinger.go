package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// BollingerBands implements Bollinger Bands (upper, middle, lower).
type BollingerBands struct {
	out    []float64
	period int
	stdDev float64
	buffer []float64
	sum    float64
	index  int
	count  int
}

func NewBollingerBands(period int, stdDev float64) *BollingerBands {
	return &BollingerBands{
		period: period,
		stdDev: stdDev,
		buffer: make([]float64, period),
		out:    make([]float64, 3),
	}
}

func (b *BollingerBands) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	upper := make([]float64, len(candles))
	middle := make([]float64, len(candles))
	lower := make([]float64, len(candles))
	b.Reset()

	for i, c := range candles {
		values := b.UpdateAll(c)
		upper[i] = values[0]
		middle[i] = values[1]
		lower[i] = values[2]
	}
	return [][]float64{upper, middle, lower}
}

func (b *BollingerBands) UpdateAll(candle indicators.OHLCV) []float64 {
	old := b.buffer[b.index]
	b.buffer[b.index] = candle.Close
	b.sum += candle.Close - old
	b.index = (b.index + 1) % b.period
	if b.count < b.period {
		b.count++
	}
	if b.count < b.period {
		b.out[0] = math.NaN()
		b.out[1] = math.NaN()
		b.out[2] = math.NaN()
		return b.out
	}

	sma := b.sum / float64(b.period)

	// Standard deviation
	variance := 0.0
	for i := 0; i < b.period; i++ {
		diff := b.buffer[i] - sma
		variance += diff * diff
	}
	sd := math.Sqrt(variance / float64(b.period))

	upper := sma + b.stdDev*sd
	lower := sma - b.stdDev*sd
	b.out[0] = upper
	b.out[1] = sma
	b.out[2] = lower
	return b.out
}

func (b *BollingerBands) Reset() {
	b.buffer = make([]float64, b.period)
	b.sum = 0
	b.index = 0
	b.count = 0
}

func (b *BollingerBands) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Bollinger Bands",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(b.period), Min: 2, Max: 200, Step: 1},
			{Name: "Std Dev", DefaultValue: b.stdDev, Min: 0.5, Max: 5, Step: 0.5},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Upper", Color: "#EF5350", FillColor: "#EF535020", Style: indicators.StyleLine, Width: 1, Opacity: 1.0},
			{Name: "Middle", Color: "#2196F3", Style: indicators.StyleDashed, Width: 1, Opacity: 1.0, DashArray: []int{5, 5}},
			{Name: "Lower", Color: "#66BB6A", FillColor: "#66BB6A20", Style: indicators.StyleLine, Width: 1, Opacity: 1.0},
		},
	}
}
