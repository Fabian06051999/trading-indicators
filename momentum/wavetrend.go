package momentum

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// WaveTrend implements the WaveTrend oscillator (LazyBear).
type WaveTrend struct {
	out       []float64
	chPeriod  int
	avgPeriod int
	ema1      *moving_averages.EMA
	ema2      *moving_averages.EMA
	ema3      *moving_averages.EMA
	d         []float64
	dIndex    int
	dCount    int
}

func NewWaveTrend(chPeriod, avgPeriod int) *WaveTrend {
	return &WaveTrend{
		chPeriod:  chPeriod,
		avgPeriod: avgPeriod,
		ema1:      moving_averages.NewEMA(chPeriod),
		ema2:      moving_averages.NewEMA(chPeriod),
		ema3:      moving_averages.NewEMA(avgPeriod),
		d:         make([]float64, 4),
		out:       make([]float64, 2),
	}
}

func (w *WaveTrend) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	wt1Values := make([]float64, len(candles))
	wt2Values := make([]float64, len(candles))
	w.Reset()

	for i, c := range candles {
		values := w.UpdateAll(c)
		wt1Values[i] = values[0]
		wt2Values[i] = values[1]
	}
	return [][]float64{wt1Values, wt2Values}
}

func (w *WaveTrend) UpdateAll(candle indicators.OHLCV) []float64 {
	hlc3 := (candle.High + candle.Low + candle.Close) / 3.0

	esa := w.ema1.UpdateAll(indicators.OHLCV{Close: hlc3})[0]
	if math.IsNaN(esa) {
		w.out[0] = math.NaN()
		w.out[1] = math.NaN()
		return w.out
	}

	d := math.Abs(hlc3 - esa)
	dd := w.ema2.UpdateAll(indicators.OHLCV{Close: d})[0]
	if math.IsNaN(dd) {
		w.out[0] = math.NaN()
		w.out[1] = math.NaN()
		return w.out
	}

	ci := 0.0
	if dd != 0 {
		ci = (hlc3 - esa) / (0.015 * dd)
	}

	wt1 := w.ema3.UpdateAll(indicators.OHLCV{Close: ci})[0]
	if math.IsNaN(wt1) {
		w.out[0] = math.NaN()
		w.out[1] = math.NaN()
		return w.out
	}

	// WT2 = SMA(WT1, 4)
	w.d[w.dIndex] = wt1
	w.dIndex = (w.dIndex + 1) % 4
	w.dCount++

	if w.dCount < 4 {
		w.out[0] = wt1
		w.out[1] = 0
		return w.out
	}

	wt2 := (w.d[0] + w.d[1] + w.d[2] + w.d[3]) / 4.0
	w.out[0] = wt1
	w.out[1] = wt2
	return w.out
}

func (w *WaveTrend) Reset() {
	w.ema1.Reset()
	w.ema2.Reset()
	w.ema3.Reset()
	w.d = make([]float64, 4)
	w.dIndex = 0
	w.dCount = 0
}

func (w *WaveTrend) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "WaveTrend",
		Parameters: []indicators.Parameter{
			{Name: "Channel Period", DefaultValue: float64(w.chPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "Average Period", DefaultValue: float64(w.avgPeriod), Min: 2, Max: 50, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "WT1",
				Color: "#00C853",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 60, Label: "Overbought", Color: "#EF5350"},
					{Value: -60, Label: "Oversold", Color: "#66BB6A"},
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
			{Name: "WT2", Color: "#D50000", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
