package moving_averages

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// McGinleyDynamic implements the McGinley Dynamic indicator.
type McGinleyDynamic struct {
	period int
	value  float64
	count  int
}

func NewMcGinleyDynamic(period int) *McGinleyDynamic {
	return &McGinleyDynamic{
		period: period,
	}
}

func (m *McGinleyDynamic) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	m.Reset()

	for i, c := range candles {
		result[i] = m.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (m *McGinleyDynamic) UpdateAll(candle indicators.OHLCV) []float64 {
	m.count++
	if m.count == 1 {
		m.value = candle.Close
		return []float64{m.value}
	}

	if m.value == 0 {
		m.value = candle.Close
		return []float64{m.value}
	}

	ratio := candle.Close / m.value
	denom := float64(m.period) * math.Pow(ratio, 4)
	if denom == 0 {
		return []float64{m.value}
	}

	m.value = m.value + (candle.Close-m.value)/denom
	return []float64{m.value}
}

func (m *McGinleyDynamic) Reset() {
	m.value = 0
	m.count = 0
}

func (m *McGinleyDynamic) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "McGinley Dynamic",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(m.period), Min: 2, Max: 500, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "McGinley", Color: "#00897B", Style: indicators.StyleLine, Width: 2},
		},
	}
}
