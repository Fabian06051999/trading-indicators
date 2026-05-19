package oscillators

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// MFI implements the Money Flow Index.
type MFI struct {
	period     int
	posFlows   []float64
	negFlows   []float64
	prevTP     float64
	index      int
	count      int
}

func NewMFI(period int) *MFI {
	period = indicators.ClampMin(period, 2)
	return &MFI{
		period:   period,
		posFlows: make([]float64, period),
		negFlows: make([]float64, period),
	}
}

func (m *MFI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	m.Reset()

	for i, c := range candles {
		result[i] = m.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (m *MFI) UpdateAll(candle indicators.OHLCV) []float64 {
	tp := (candle.High + candle.Low + candle.Close) / 3.0
	mf := tp * candle.Volume
	m.count++

	if m.count == 1 {
		m.prevTP = tp
		return []float64{math.NaN()}
	}

	// Classify money flow
	m.posFlows[m.index] = 0
	m.negFlows[m.index] = 0
	if tp > m.prevTP {
		m.posFlows[m.index] = mf
	} else if tp < m.prevTP {
		m.negFlows[m.index] = mf
	}
	m.prevTP = tp
	m.index = (m.index + 1) % m.period

	if m.count <= m.period {
		return []float64{math.NaN()}
	}

	posMF := 0.0
	negMF := 0.0
	for i := 0; i < m.period; i++ {
		posMF += m.posFlows[i]
		negMF += m.negFlows[i]
	}

	if negMF == 0 {
		return []float64{100}
	}
	mfRatio := posMF / negMF
	return []float64{100 - 100/(1+mfRatio)}
}

func (m *MFI) Reset() {
	m.posFlows = make([]float64, m.period)
	m.negFlows = make([]float64, m.period)
	m.prevTP = 0
	m.index = 0
	m.count = 0
}

func (m *MFI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Money Flow Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(m.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "MFI",
				Color:  "#26A69A",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 80, Label: "Overbought", Color: "#EF5350"},
					{Value: 20, Label: "Oversold", Color: "#66BB6A"},
				},
			},
		},
	}
}
