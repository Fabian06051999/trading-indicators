package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// ElderRay implements Elder's Bull/Bear Power.
type ElderRay struct {
	period int
	ema    *moving_averages.EMA
}

func NewElderRay(period int) *ElderRay {
	return &ElderRay{
		period: period,
		ema:    moving_averages.NewEMA(period),
	}
}

func (e *ElderRay) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	bull := make([]float64, len(candles))
	bear := make([]float64, len(candles))
	e.Reset()

	for i, c := range candles {
		values := e.UpdateAll(c)
		bull[i] = values[0]
		bear[i] = values[1]
	}
	return [][]float64{bull, bear}
}

func (e *ElderRay) UpdateAll(candle indicators.OHLCV) []float64 {
	emaVal := e.ema.UpdateAll(candle)[0]
	if emaVal == 0 {
		return []float64{math.NaN(), math.NaN()}
	}

	bullPower := candle.High - emaVal
	bearPower := candle.Low - emaVal
	return []float64{bullPower, bearPower}
}

func (e *ElderRay) Reset() {
	e.ema.Reset()
}

func (e *ElderRay) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Elder Ray (Bull/Bear Power)",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(e.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "Bull Power",
				Color: "#4CAF50",
				Style: indicators.StyleHistogram,
				Width: 1,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
			{Name: "Bear Power", Color: "#F44336", Style: indicators.StyleHistogram, Width: 1},
		},
	}
}
