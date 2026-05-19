package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// ROC implements the Rate of Change indicator.
type ROC struct {
	period int
	buffer []float64
	index  int
	count  int
}

func NewROC(period int) *ROC {
	period = indicators.ClampMin(period, 1)
	return &ROC{
		period: period,
		buffer: make([]float64, period+1),
	}
}

func (r *ROC) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	r.Reset()

	for i, c := range candles {
		result[i] = r.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (r *ROC) UpdateAll(candle indicators.OHLCV) []float64 {
	r.buffer[r.index] = candle.Close
	r.count++

	if r.count <= r.period {
		r.index = (r.index + 1) % (r.period + 1)
		return []float64{math.NaN()}
	}

	pastIdx := (r.index + 1) % (r.period + 1)
	pastVal := r.buffer[pastIdx]
	r.index = (r.index + 1) % (r.period + 1)

	if pastVal == 0 {
		return []float64{math.NaN()}
	}
	return []float64{((candle.Close - pastVal) / pastVal) * 100}
}

func (r *ROC) Reset() {
	r.buffer = make([]float64, r.period+1)
	r.index = 0
	r.count = 0
}

func (r *ROC) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Rate of Change",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(r.period), Min: 1, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "ROC",
				Color: "#E91E63",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
		},
	}
}
