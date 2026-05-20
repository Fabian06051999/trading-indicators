package trend

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
	"math"
)

// MACD implements the Moving Average Convergence Divergence.
type MACD struct {
	out          []float64
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
	fastEMA      *moving_averages.EMA
	slowEMA      *moving_averages.EMA
	signalEMA    *moving_averages.EMA
	count        int
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		fastPeriod:   fastPeriod,
		slowPeriod:   slowPeriod,
		signalPeriod: signalPeriod,
		fastEMA:      moving_averages.NewEMA(fastPeriod),
		slowEMA:      moving_averages.NewEMA(slowPeriod),
		signalEMA:    moving_averages.NewEMA(signalPeriod),
		out:          make([]float64, 3),
	}
}

func (m *MACD) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	macdLine := make([]float64, len(candles))
	signalLine := make([]float64, len(candles))
	histogram := make([]float64, len(candles))
	m.Reset()

	for i, c := range candles {
		values := m.UpdateAll(c)
		macdLine[i] = values[0]
		signalLine[i] = values[1]
		histogram[i] = values[2]
	}
	return [][]float64{macdLine, signalLine, histogram}
}

func (m *MACD) UpdateAll(candle indicators.OHLCV) []float64 {
	m.count++
	fast := m.fastEMA.UpdateAll(candle)[0]
	slow := m.slowEMA.UpdateAll(candle)[0]

	if math.IsNaN(fast) || math.IsNaN(slow) {
		m.out[0] = math.NaN()
		m.out[1] = math.NaN()
		m.out[2] = math.NaN()
		return m.out
	}

	macdVal := fast - slow
	signal := m.signalEMA.UpdateAll(indicators.OHLCV{Close: macdVal})[0]

	if math.IsNaN(signal) {
		m.out[0] = macdVal
		m.out[1] = math.NaN()
		m.out[2] = math.NaN()
		return m.out
	}

	hist := macdVal - signal
	m.out[0] = macdVal
	m.out[1] = signal
	m.out[2] = hist
	return m.out
}

func (m *MACD) Reset() {
	m.fastEMA.Reset()
	m.slowEMA.Reset()
	m.signalEMA.Reset()
	m.count = 0
}

func (m *MACD) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "MACD",
		Parameters: []indicators.Parameter{
			{Name: "Fast Period", DefaultValue: float64(m.fastPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Slow Period", DefaultValue: float64(m.slowPeriod), Min: 2, Max: 200, Step: 1},
			{Name: "Signal Period", DefaultValue: float64(m.signalPeriod), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "MACD", Color: "#2196F3", Style: indicators.StyleLine, Width: 2, Opacity: 1.0},
			{Name: "Signal", Color: "#FF9800", Style: indicators.StyleLine, Width: 1, Opacity: 1.0},
			{Name: "Histogram", Color: "#4CAF50", UpColor: "#4CAF50", DownColor: "#F44336", Style: indicators.StyleHistogram, Width: 1, Opacity: 0.8},
		},
	}
}
