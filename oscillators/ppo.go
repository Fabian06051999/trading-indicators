package oscillators

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// PPO implements the Percentage Price Oscillator.
type PPO struct {
	fastPeriod   int
	slowPeriod   int
	signalPeriod int
	fastEMA      *moving_averages.EMA
	slowEMA      *moving_averages.EMA
	signalEMA    *moving_averages.EMA
	count        int
}

func NewPPO(fastPeriod, slowPeriod, signalPeriod int) *PPO {
	return &PPO{
		fastPeriod:   fastPeriod,
		slowPeriod:   slowPeriod,
		signalPeriod: signalPeriod,
		fastEMA:      moving_averages.NewEMA(fastPeriod),
		slowEMA:      moving_averages.NewEMA(slowPeriod),
		signalEMA:    moving_averages.NewEMA(signalPeriod),
	}
}

func (p *PPO) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	ppoLine := make([]float64, len(candles))
	signalLine := make([]float64, len(candles))
	histogram := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		values := p.UpdateAll(c)
		ppoLine[i] = values[0]
		signalLine[i] = values[1]
		histogram[i] = values[2]
	}
	return [][]float64{ppoLine, signalLine, histogram}
}

func (p *PPO) UpdateAll(candle indicators.OHLCV) []float64 {
	p.count++
	fast := p.fastEMA.Update(candle)
	slow := p.slowEMA.Update(candle)

	if fast == 0 || slow == 0 {
		return []float64{math.NaN(), math.NaN(), math.NaN()}
	}

	ppoVal := ((fast - slow) / slow) * 100
	signal := p.signalEMA.Update(indicators.OHLCV{Close: ppoVal})

	if signal == 0 {
		return []float64{ppoVal, 0, 0}
	}

	return []float64{ppoVal, signal, ppoVal - signal}
}

func (p *PPO) Reset() {
	p.fastEMA.Reset()
	p.slowEMA.Reset()
	p.signalEMA.Reset()
	p.count = 0
}

func (p *PPO) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Percentage Price Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Fast Period", DefaultValue: float64(p.fastPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Slow Period", DefaultValue: float64(p.slowPeriod), Min: 2, Max: 200, Step: 1},
			{Name: "Signal Period", DefaultValue: float64(p.signalPeriod), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "PPO", Color: "#2196F3", Style: indicators.StyleLine, Width: 2},
			{Name: "Signal", Color: "#FF9800", Style: indicators.StyleLine, Width: 1},
			{Name: "Histogram", Color: "#4CAF50", Style: indicators.StyleHistogram, Width: 1},
		},
	}
}
