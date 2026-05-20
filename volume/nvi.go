package volume

import (
	"github.com/Fabian06051999/trading-indicators"
)

// NVI implements the Negative Volume Index.
type NVI struct {
	out       []float64
	value     float64
	prevClose float64
	prevVol   float64
	count     int
}

func NewNVI() *NVI {
	return &NVI{value: 1000, out: make([]float64, 1)}
}

func (n *NVI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	n.Reset()

	for i, c := range candles {
		result[i] = n.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (n *NVI) UpdateAll(candle indicators.OHLCV) []float64 {
	n.count++

	if n.count == 1 {
		n.prevClose = candle.Close
		n.prevVol = candle.Volume
		n.out[0] = n.value
		return n.out
	}

	if candle.Volume < n.prevVol && n.prevClose != 0 {
		roc := (candle.Close - n.prevClose) / n.prevClose
		n.value += n.value * roc
	}

	n.prevClose = candle.Close
	n.prevVol = candle.Volume
	n.out[0] = n.value
	return n.out
}

func (n *NVI) Reset() {
	n.value = 1000
	n.prevClose = 0
	n.prevVol = 0
	n.count = 0
}

func (n *NVI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name:       "Negative Volume Index",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "NVI", Color: "#5D4037", Style: indicators.StyleLine, Width: 2},
		},
	}
}

// PVI implements the Positive Volume Index.
type PVI struct {
	out       []float64
	value     float64
	prevClose float64
	prevVol   float64
	count     int
}

func NewPVI() *PVI {
	return &PVI{value: 1000, out: make([]float64, 1)}
}

func (p *PVI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		result[i] = p.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (p *PVI) UpdateAll(candle indicators.OHLCV) []float64 {
	p.count++

	if p.count == 1 {
		p.prevClose = candle.Close
		p.prevVol = candle.Volume
		p.out[0] = p.value
		return p.out
	}

	if candle.Volume > p.prevVol && p.prevClose != 0 {
		roc := (candle.Close - p.prevClose) / p.prevClose
		p.value += p.value * roc
	}

	p.prevClose = candle.Close
	p.prevVol = candle.Volume
	p.out[0] = p.value
	return p.out
}

func (p *PVI) Reset() {
	p.value = 1000
	p.prevClose = 0
	p.prevVol = 0
	p.count = 0
}

func (p *PVI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name:       "Positive Volume Index",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "PVI", Color: "#2E7D32", Style: indicators.StyleLine, Width: 2},
		},
	}
}
