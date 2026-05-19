package volume

import (
	"github.com/Fabian06051999/trading-indicators"
)

// OBV implements On Balance Volume.
type OBV struct {
	value     float64
	prevClose float64
	count     int
}

func NewOBV() *OBV {
	return &OBV{}
}

func (o *OBV) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	o.Reset()

	for i, c := range candles {
		result[i] = o.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (o *OBV) UpdateAll(candle indicators.OHLCV) []float64 {
	o.count++
	if o.count == 1 {
		o.prevClose = candle.Close
		o.value = candle.Volume
		return []float64{o.value}
	}

	if candle.Close > o.prevClose {
		o.value += candle.Volume
	} else if candle.Close < o.prevClose {
		o.value -= candle.Volume
	}
	o.prevClose = candle.Close
	return []float64{o.value}
}

func (o *OBV) Reset() {
	o.value = 0
	o.prevClose = 0
	o.count = 0
}

func (o *OBV) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "On Balance Volume",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "OBV", Color: "#3F51B5", Style: indicators.StyleLine, Width: 2},
		},
	}
}
