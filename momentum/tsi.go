package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// TSI implements the True Strength Index.
type TSI struct {
	longPeriod  int
	shortPeriod int
	signalPeriod int
	longEMA1    *moving_averages.EMA
	shortEMA1   *moving_averages.EMA
	longEMA2    *moving_averages.EMA
	shortEMA2   *moving_averages.EMA
	signalEMA   *moving_averages.EMA
	prevClose   float64
	count       int
}

func NewTSI(longPeriod, shortPeriod, signalPeriod int) *TSI {
	return &TSI{
		longPeriod:   longPeriod,
		shortPeriod:  shortPeriod,
		signalPeriod: signalPeriod,
		longEMA1:     moving_averages.NewEMA(longPeriod),
		shortEMA1:    moving_averages.NewEMA(shortPeriod),
		longEMA2:     moving_averages.NewEMA(longPeriod),
		shortEMA2:    moving_averages.NewEMA(shortPeriod),
		signalEMA:    moving_averages.NewEMA(signalPeriod),
	}
}

func (t *TSI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	tsiValues := make([]float64, len(candles))
	signalValues := make([]float64, len(candles))
	t.Reset()

	for i, c := range candles {
		values := t.UpdateAll(c)
		tsiValues[i] = values[0]
		signalValues[i] = values[1]
	}
	return [][]float64{tsiValues, signalValues}
}

func (t *TSI) UpdateAll(candle indicators.OHLCV) []float64 {
	t.count++

	if t.count == 1 {
		t.prevClose = candle.Close
		return []float64{math.NaN(), math.NaN()}
	}

	pc := candle.Close - t.prevClose
	t.prevClose = candle.Close

	// Double-smoothed price change
	ds1 := t.longEMA1.Update(indicators.OHLCV{Close: pc})
	if ds1 == 0 {
		return []float64{math.NaN(), math.NaN()}
	}
	ds := t.shortEMA1.Update(indicators.OHLCV{Close: ds1})
	if ds == 0 {
		return []float64{math.NaN(), math.NaN()}
	}

	// Double-smoothed absolute price change
	apc := pc
	if apc < 0 {
		apc = -apc
	}
	ads1 := t.longEMA2.Update(indicators.OHLCV{Close: apc})
	if ads1 == 0 {
		return []float64{math.NaN(), math.NaN()}
	}
	ads := t.shortEMA2.Update(indicators.OHLCV{Close: ads1})
	if ads == 0 {
		return []float64{math.NaN(), math.NaN()}
	}

	tsiVal := 0.0
	if ads != 0 {
		tsiVal = (ds / ads) * 100
	}

	signal := t.signalEMA.Update(indicators.OHLCV{Close: tsiVal})
	return []float64{tsiVal, signal}
}

func (t *TSI) Reset() {
	t.longEMA1.Reset()
	t.shortEMA1.Reset()
	t.longEMA2.Reset()
	t.shortEMA2.Reset()
	t.signalEMA.Reset()
	t.prevClose = 0
	t.count = 0
}

func (t *TSI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "True Strength Index",
		Parameters: []indicators.Parameter{
			{Name: "Long Period", DefaultValue: float64(t.longPeriod), Min: 5, Max: 50, Step: 1},
			{Name: "Short Period", DefaultValue: float64(t.shortPeriod), Min: 2, Max: 30, Step: 1},
			{Name: "Signal Period", DefaultValue: float64(t.signalPeriod), Min: 2, Max: 30, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "TSI",
				Color: "#1976D2",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
					{Value: 25, Label: "OB", Color: "#EF5350"},
					{Value: -25, Label: "OS", Color: "#66BB6A"},
				},
			},
			{Name: "Signal", Color: "#FF9800", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
