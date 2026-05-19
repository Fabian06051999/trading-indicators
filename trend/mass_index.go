package trend

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// MassIndex implements the Mass Index (trend reversal detection).
type MassIndex struct {
	emaPeriod  int
	sumPeriod  int
	ema1       *moving_averages.EMA
	ema2       *moving_averages.EMA
	ratios     []float64
	index      int
	count      int
}

func NewMassIndex(emaPeriod, sumPeriod int) *MassIndex {
	return &MassIndex{
		emaPeriod: emaPeriod,
		sumPeriod: sumPeriod,
		ema1:      moving_averages.NewEMA(emaPeriod),
		ema2:      moving_averages.NewEMA(emaPeriod),
		ratios:    make([]float64, sumPeriod),
	}
}

func (m *MassIndex) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	m.Reset()

	for i, c := range candles {
		result[i] = m.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (m *MassIndex) UpdateAll(candle indicators.OHLCV) []float64 {
	hl := candle.High - candle.Low
	e1 := m.ema1.UpdateAll(indicators.OHLCV{Close: hl})[0]
	if e1 == 0 {
		return []float64{math.NaN()}
	}
	e2 := m.ema2.UpdateAll(indicators.OHLCV{Close: e1})[0]
	if e2 == 0 {
		return []float64{math.NaN()}
	}

	ratio := 0.0
	if e2 != 0 {
		ratio = e1 / e2
	}

	m.ratios[m.index] = ratio
	m.index = (m.index + 1) % m.sumPeriod
	m.count++

	if m.count < m.sumPeriod {
		return []float64{math.NaN()}
	}

	sum := 0.0
	for i := 0; i < m.sumPeriod; i++ {
		sum += m.ratios[i]
	}
	return []float64{sum}
}

func (m *MassIndex) Reset() {
	m.ema1.Reset()
	m.ema2.Reset()
	m.ratios = make([]float64, m.sumPeriod)
	m.index = 0
	m.count = 0
}

func (m *MassIndex) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Mass Index",
		Parameters: []indicators.Parameter{
			{Name: "EMA Period", DefaultValue: float64(m.emaPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "Sum Period", DefaultValue: float64(m.sumPeriod), Min: 5, Max: 50, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Mass Index",
				Color: "#AD1457",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 27, Label: "Bulge", Color: "#EF5350"},
					{Value: 26.5, Label: "Trigger", Color: "#66BB6A"},
				},
			},
		},
	}
}
