package volume

import (
	"github.com/Fabian06051999/trading-indicators"
)

// NVI implements the Negative Volume Index.
type NVI struct {
	value     float64
	prevClose float64
	prevVol   float64
	count     int
}

func NewNVI() *NVI {
	return &NVI{value: 1000}
}

func (n *NVI) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	n.Reset()

	for i, c := range candles {
		result[i] = n.Update(c)
	}
	return result
}

func (n *NVI) Update(candle indicators.OHLCV) float64 {
	n.count++

	if n.count == 1 {
		n.prevClose = candle.Close
		n.prevVol = candle.Volume
		return n.value
	}

	if candle.Volume < n.prevVol && n.prevClose != 0 {
		roc := (candle.Close - n.prevClose) / n.prevClose
		n.value += n.value * roc
	}

	n.prevClose = candle.Close
	n.prevVol = candle.Volume
	return n.value
}

func (n *NVI) Reset() {
	n.value = 1000
	n.prevClose = 0
	n.prevVol = 0
	n.count = 0
}

func (n *NVI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Negative Volume Index",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "NVI", Color: "#5D4037", Style: indicators.StyleLine, Width: 2},
		},
	}
}

// PVI implements the Positive Volume Index.
type PVI struct {
	value     float64
	prevClose float64
	prevVol   float64
	count     int
}

func NewPVI() *PVI {
	return &PVI{value: 1000}
}

func (p *PVI) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		result[i] = p.Update(c)
	}
	return result
}

func (p *PVI) Update(candle indicators.OHLCV) float64 {
	p.count++

	if p.count == 1 {
		p.prevClose = candle.Close
		p.prevVol = candle.Volume
		return p.value
	}

	if candle.Volume > p.prevVol && p.prevClose != 0 {
		roc := (candle.Close - p.prevClose) / p.prevClose
		p.value += p.value * roc
	}

	p.prevClose = candle.Close
	p.prevVol = candle.Volume
	return p.value
}

func (p *PVI) Reset() {
	p.value = 1000
	p.prevClose = 0
	p.prevVol = 0
	p.count = 0
}

func (p *PVI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Positive Volume Index",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "PVI", Color: "#2E7D32", Style: indicators.StyleLine, Width: 2},
		},
	}
}
