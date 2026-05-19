package oscillators

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// RSI implements the Relative Strength Index.
type RSI struct {
	period   int
	avgGain  float64
	avgLoss  float64
	prevClose float64
	count    int
	gains    float64
	losses   float64
}

func NewRSI(period int) *RSI {
	return &RSI{
		period: period,
	}
}

func (r *RSI) Calculate(candles []indicators.OHLCV) []float64 {
	result := make([]float64, len(candles))
	r.Reset()

	for i, c := range candles {
		result[i] = r.Update(c)
	}
	return result
}

func (r *RSI) Update(candle indicators.OHLCV) float64 {
	r.count++

	if r.count == 1 {
		r.prevClose = candle.Close
		return math.NaN()
	}

	change := candle.Close - r.prevClose
	r.prevClose = candle.Close

	gain := 0.0
	loss := 0.0
	if change > 0 {
		gain = change
	} else {
		loss = -change
	}

	if r.count <= r.period+1 {
		r.gains += gain
		r.losses += loss

		if r.count == r.period+1 {
			r.avgGain = r.gains / float64(r.period)
			r.avgLoss = r.losses / float64(r.period)

			if r.avgLoss == 0 {
				return 100
			}
			rs := r.avgGain / r.avgLoss
			return 100 - 100/(1+rs)
		}
		return math.NaN()
	}

	r.avgGain = (r.avgGain*float64(r.period-1) + gain) / float64(r.period)
	r.avgLoss = (r.avgLoss*float64(r.period-1) + loss) / float64(r.period)

	if r.avgLoss == 0 {
		return 100
	}
	rs := r.avgGain / r.avgLoss
	return 100 - 100/(1+rs)
}

func (r *RSI) Reset() {
	r.avgGain = 0
	r.avgLoss = 0
	r.prevClose = 0
	r.count = 0
	r.gains = 0
	r.losses = 0
}

func (r *RSI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Relative Strength Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(r.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "RSI",
				Color:  "#7B1FA2",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 70, Label: "Overbought", Color: "#EF5350"},
					{Value: 30, Label: "Oversold", Color: "#66BB6A"},
					{Value: 50, Label: "Mid", Color: "#9E9E9E"},
				},
			},
		},
	}
}
