package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/volatility"
)

// PercentB implements %B (position within Bollinger Bands).
type PercentB struct {
	bb *volatility.BollingerBands
}

func NewPercentB(period int, stdDev float64) *PercentB {
	return &PercentB{
		bb: volatility.NewBollingerBands(period, stdDev),
	}
}

func (p *PercentB) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		result[i] = p.Update(c)
	}
	return result
}

func (p *PercentB) Update(candle indicators.OHLCV) float64 {
	bands := p.bb.UpdateAll(candle)
	upper := bands[0]
	lower := bands[2]

	if upper == 0 || lower == 0 || upper == lower {
		return 0
	}

	return (candle.Close - lower) / (upper - lower)
}

func (p *PercentB) Reset() {
	p.bb.Reset()
}

func (p *PercentB) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Bollinger %B",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: 20, Min: 2, Max: 200, Step: 1},
			{Name: "Std Dev", DefaultValue: 2, Min: 0.5, Max: 5, Step: 0.5},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "%B",
				Color:  "#7B1FA2",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: -0.5, Max: 1.5},
				Levels: []indicators.Level{
					{Value: 1, Label: "Upper Band", Color: "#EF5350"},
					{Value: 0.5, Label: "Middle", Color: "#9E9E9E"},
					{Value: 0, Label: "Lower Band", Color: "#66BB6A"},
				},
			},
		},
	}
}
