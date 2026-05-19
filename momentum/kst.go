package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/moving_averages"
)

// KST implements the Know Sure Thing indicator.
type KST struct {
	roc1Period int
	roc2Period int
	roc3Period int
	roc4Period int
	sma1      *moving_averages.SMA
	sma2      *moving_averages.SMA
	sma3      *moving_averages.SMA
	sma4      *moving_averages.SMA
	signalSMA *moving_averages.SMA
	buffer    []float64
	index     int
	count     int
	maxPeriod int
}

func NewKST(roc1, roc2, roc3, roc4, sma1, sma2, sma3, sma4, signal int) *KST {
	maxP := roc4
	if roc3 > maxP {
		maxP = roc3
	}
	if roc2 > maxP {
		maxP = roc2
	}
	return &KST{
		roc1Period: roc1,
		roc2Period: roc2,
		roc3Period: roc3,
		roc4Period: roc4,
		sma1:      moving_averages.NewSMA(sma1),
		sma2:      moving_averages.NewSMA(sma2),
		sma3:      moving_averages.NewSMA(sma3),
		sma4:      moving_averages.NewSMA(sma4),
		signalSMA: moving_averages.NewSMA(signal),
		buffer:    make([]float64, maxP+1),
		maxPeriod: maxP,
	}
}

func (k *KST) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	kstValues := make([]float64, len(candles))
	signalValues := make([]float64, len(candles))
	k.Reset()

	for i, c := range candles {
		values := k.UpdateAll(c)
		kstValues[i] = values[0]
		signalValues[i] = values[1]
	}
	return [][]float64{kstValues, signalValues}
}

func (k *KST) UpdateAll(candle indicators.OHLCV) []float64 {
	k.buffer[k.index] = candle.Close
	k.count++

	if k.count <= k.maxPeriod {
		k.index = (k.index + 1) % (k.maxPeriod + 1)
		return []float64{0, 0}
	}

	// ROC calculations
	roc1 := k.roc(k.roc1Period, candle.Close)
	roc2 := k.roc(k.roc2Period, candle.Close)
	roc3 := k.roc(k.roc3Period, candle.Close)
	roc4 := k.roc(k.roc4Period, candle.Close)

	k.index = (k.index + 1) % (k.maxPeriod + 1)

	// Smooth each ROC with SMA
	s1 := k.sma1.Update(indicators.OHLCV{Close: roc1})
	s2 := k.sma2.Update(indicators.OHLCV{Close: roc2})
	s3 := k.sma3.Update(indicators.OHLCV{Close: roc3})
	s4 := k.sma4.Update(indicators.OHLCV{Close: roc4})

	if s1 == 0 || s2 == 0 || s3 == 0 || s4 == 0 {
		return []float64{0, 0}
	}

	kstVal := s1*1 + s2*2 + s3*3 + s4*4
	signal := k.signalSMA.Update(indicators.OHLCV{Close: kstVal})

	return []float64{kstVal, signal}
}

func (k *KST) roc(period int, current float64) float64 {
	pastIdx := (k.index - period + k.maxPeriod + 1) % (k.maxPeriod + 1)
	past := k.buffer[pastIdx]
	if past == 0 {
		return 0
	}
	return ((current - past) / past) * 100
}

func (k *KST) Reset() {
	k.sma1.Reset()
	k.sma2.Reset()
	k.sma3.Reset()
	k.sma4.Reset()
	k.signalSMA.Reset()
	k.buffer = make([]float64, k.maxPeriod+1)
	k.index = 0
	k.count = 0
}

func (k *KST) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Know Sure Thing",
		Parameters: []indicators.Parameter{
			{Name: "ROC1", DefaultValue: float64(k.roc1Period), Min: 5, Max: 30, Step: 1},
			{Name: "ROC2", DefaultValue: float64(k.roc2Period), Min: 5, Max: 30, Step: 1},
			{Name: "ROC3", DefaultValue: float64(k.roc3Period), Min: 5, Max: 50, Step: 1},
			{Name: "ROC4", DefaultValue: float64(k.roc4Period), Min: 10, Max: 50, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:  "KST",
				Color: "#0D47A1",
				Style: indicators.StyleLine,
				Width: 2,
				Levels: []indicators.Level{
					{Value: 0, Label: "Zero", Color: "#9E9E9E"},
				},
			},
			{Name: "Signal", Color: "#F44336", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
