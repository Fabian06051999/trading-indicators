package oscillators

import (
	"github.com/Fabian06051999/trading-indicators"
)

// AcceleratorOscillator implements Bill Williams' Accelerator Oscillator.
// AC = AO - SMA(AO, 5)
type AcceleratorOscillator struct {
	ao        *AwesomeOscillator
	aoBuffer  []float64
	aoSum     float64
	aoIndex   int
	aoCount   int
}

func NewAcceleratorOscillator() *AcceleratorOscillator {
	return &AcceleratorOscillator{
		ao:       NewAwesomeOscillator(),
		aoBuffer: make([]float64, 5),
	}
}

func (ac *AcceleratorOscillator) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	ac.Reset()

	for i, c := range candles {
		result[i] = ac.Update(c)
	}
	return result
}

func (ac *AcceleratorOscillator) Update(candle indicators.OHLCV) float64 {
	aoVal := ac.ao.Update(candle)
	if aoVal == 0 {
		return 0
	}

	ac.aoCount++
	ac.aoSum -= ac.aoBuffer[ac.aoIndex]
	ac.aoBuffer[ac.aoIndex] = aoVal
	ac.aoSum += aoVal
	ac.aoIndex = (ac.aoIndex + 1) % 5

	if ac.aoCount < 5 {
		return 0
	}

	aoSMA := ac.aoSum / 5.0
	return aoVal - aoSMA
}

func (ac *AcceleratorOscillator) Reset() {
	ac.ao.Reset()
	ac.aoBuffer = make([]float64, 5)
	ac.aoSum = 0
	ac.aoIndex = 0
	ac.aoCount = 0
}

func (ac *AcceleratorOscillator) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Accelerator Oscillator",
		Parameters: []indicators.Parameter{},
		Pane:       indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "AC",
				Color: "#8D6E63",
				Style: indicators.StyleHistogram,
				Width: 1,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
