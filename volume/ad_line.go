package volume

import (
	"github.com/Fabian06051999/trading-indicators"
)

// ADLine implements the Accumulation/Distribution Line.
type ADLine struct {
	value float64
}

func NewADLine() *ADLine {
	return &ADLine{}
}

func (a *ADLine) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		result[i] = a.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (a *ADLine) UpdateAll(candle indicators.OHLCV) []float64 {
	hl := candle.High - candle.Low
	if hl == 0 {
		return []float64{a.value}
	}
	mfm := ((candle.Close - candle.Low) - (candle.High - candle.Close)) / hl
	a.value += mfm * candle.Volume
	return []float64{a.value}
}

func (a *ADLine) Reset() {
	a.value = 0
}

func (a *ADLine) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Accumulation/Distribution Line",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "A/D", Color: "#1565C0", Style: indicators.StyleLine, Width: 2},
		},
	}
}
