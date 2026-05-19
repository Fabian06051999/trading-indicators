package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// KeltnerChannel implements Keltner Channels.
type KeltnerChannel struct {
	emaPeriod  int
	atrPeriod  int
	multiplier float64
	ema        *moving_averages.EMA
	atr        *ATR
}

func NewKeltnerChannel(emaPeriod, atrPeriod int, multiplier float64) *KeltnerChannel {
	return &KeltnerChannel{
		emaPeriod:  emaPeriod,
		atrPeriod:  atrPeriod,
		multiplier: multiplier,
		ema:        moving_averages.NewEMA(emaPeriod),
		atr:        NewATR(atrPeriod),
	}
}

func (k *KeltnerChannel) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	upper := make([]float64, len(candles))
	middle := make([]float64, len(candles))
	lower := make([]float64, len(candles))
	k.Reset()

	for i, c := range candles {
		values := k.UpdateAll(c)
		upper[i] = values[0]
		middle[i] = values[1]
		lower[i] = values[2]
	}
	return [][]float64{upper, middle, lower}
}

func (k *KeltnerChannel) UpdateAll(candle indicators.OHLCV) []float64 {
	mid := k.ema.Update(candle)
	atrVal := k.atr.Update(candle)
	_ = math.Abs(0) // keep import

	if mid == 0 || atrVal == 0 {
		return []float64{math.NaN(), math.NaN(), math.NaN()}
	}

	offset := k.multiplier * atrVal
	return []float64{mid + offset, mid, mid - offset}
}

func (k *KeltnerChannel) Reset() {
	k.ema.Reset()
	k.atr.Reset()
}

func (k *KeltnerChannel) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Keltner Channel",
		Parameters: []indicators.Parameter{
			{Name: "EMA Period", DefaultValue: float64(k.emaPeriod), Min: 2, Max: 200, Step: 1},
			{Name: "ATR Period", DefaultValue: float64(k.atrPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Multiplier", DefaultValue: k.multiplier, Min: 0.5, Max: 5, Step: 0.5},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Upper", Color: "#EF5350", Style: indicators.StyleLine, Width: 1},
			{Name: "Middle", Color: "#2196F3", Style: indicators.StyleDashed, Width: 1},
			{Name: "Lower", Color: "#66BB6A", Style: indicators.StyleLine, Width: 1},
		},
	}
}
