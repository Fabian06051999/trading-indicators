package oscillators

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// AwesomeOscillator implements Bill Williams' Awesome Oscillator.
// AO = SMA(median, 5) - SMA(median, 34)
type AwesomeOscillator struct {
	fastBuffer []float64
	slowBuffer []float64
	fastSum    float64
	slowSum    float64
	fastIndex  int
	slowIndex  int
	count      int
}

func NewAwesomeOscillator() *AwesomeOscillator {
	return &AwesomeOscillator{
		fastBuffer: make([]float64, 5),
		slowBuffer: make([]float64, 34),
	}
}

func (ao *AwesomeOscillator) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	ao.Reset()

	for i, c := range candles {
		result[i] = ao.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (ao *AwesomeOscillator) UpdateAll(candle indicators.OHLCV) []float64 {
	median := (candle.High + candle.Low) / 2.0
	ao.count++

	// Fast SMA (5)
	ao.fastSum -= ao.fastBuffer[ao.fastIndex]
	ao.fastBuffer[ao.fastIndex] = median
	ao.fastSum += median
	ao.fastIndex = (ao.fastIndex + 1) % 5

	// Slow SMA (34)
	ao.slowSum -= ao.slowBuffer[ao.slowIndex]
	ao.slowBuffer[ao.slowIndex] = median
	ao.slowSum += median
	ao.slowIndex = (ao.slowIndex + 1) % 34

	if ao.count < 34 {
		return []float64{math.NaN()}
	}

	fastSMA := ao.fastSum / 5.0
	slowSMA := ao.slowSum / 34.0
	return []float64{fastSMA - slowSMA}
}

func (ao *AwesomeOscillator) Reset() {
	ao.fastBuffer = make([]float64, 5)
	ao.slowBuffer = make([]float64, 34)
	ao.fastSum = 0
	ao.slowSum = 0
	ao.fastIndex = 0
	ao.slowIndex = 0
	ao.count = 0
}

func (ao *AwesomeOscillator) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Awesome Oscillator",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "AO",
				Color: "#4CAF50",
				Style: indicators.StyleHistogram,
				Width: 1,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
