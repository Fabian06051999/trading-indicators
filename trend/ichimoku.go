package trend

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// Ichimoku implements the Ichimoku Cloud (Kinko Hyo).
type Ichimoku struct {
	tenkanPeriod  int
	kijunPeriod   int
	senkouBPeriod int
	displacement  int
	highs         []float64
	lows          []float64
	index         int
	count         int
}

func NewIchimoku(tenkan, kijun, senkouB, displacement int) *Ichimoku {
	maxP := senkouB
	if kijun > maxP {
		maxP = kijun
	}
	return &Ichimoku{
		tenkanPeriod:  tenkan,
		kijunPeriod:   kijun,
		senkouBPeriod: senkouB,
		displacement:  displacement,
		highs:         make([]float64, maxP),
		lows:          make([]float64, maxP),
	}
}

func (ich *Ichimoku) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	tenkan := make([]float64, len(candles))
	kijun := make([]float64, len(candles))
	senkouA := make([]float64, len(candles))
	senkouB := make([]float64, len(candles))
	chikou := make([]float64, len(candles))
	ich.Reset()

	for i, c := range candles {
		values := ich.UpdateAll(c)
		tenkan[i] = values[0]
		kijun[i] = values[1]
		senkouA[i] = values[2]
		senkouB[i] = values[3]
		chikou[i] = values[4]
	}
	return [][]float64{tenkan, kijun, senkouA, senkouB, chikou}
}

func (ich *Ichimoku) UpdateAll(candle indicators.OHLCV) []float64 {
	maxP := ich.senkouBPeriod
	if ich.kijunPeriod > maxP {
		maxP = ich.kijunPeriod
	}

	ich.highs[ich.index] = candle.High
	ich.lows[ich.index] = candle.Low
	ich.count++
	ich.index = (ich.index + 1) % maxP

	// Tenkan-sen (Conversion Line)
	tenkanVal := ich.midpoint(ich.tenkanPeriod)

	// Kijun-sen (Base Line)
	kijunVal := ich.midpoint(ich.kijunPeriod)

	// Senkou Span A (Leading Span A) = (Tenkan + Kijun) / 2, displaced forward
	senkouA := 0.0
	if tenkanVal != 0 && kijunVal != 0 {
		senkouA = (tenkanVal + kijunVal) / 2.0
	}

	// Senkou Span B (Leading Span B) = midpoint of senkouBPeriod, displaced forward
	senkouB := ich.midpoint(ich.senkouBPeriod)

	// Chikou Span = current close (plotted displaced backwards)
	chikou := candle.Close

	return []float64{tenkanVal, kijunVal, senkouA, senkouB, chikou}
}

func (ich *Ichimoku) midpoint(period int) float64 {
	if ich.count < period {
		return math.NaN()
	}

	maxP := ich.senkouBPeriod
	if ich.kijunPeriod > maxP {
		maxP = ich.kijunPeriod
	}

	hh := -1e308
	ll := 1e308
	for i := 0; i < period; i++ {
		idx := (ich.index - 1 - i + maxP) % maxP
		if ich.highs[idx] > hh {
			hh = ich.highs[idx]
		}
		if ich.lows[idx] < ll {
			ll = ich.lows[idx]
		}
	}
	return (hh + ll) / 2.0
}

func (ich *Ichimoku) Reset() {
	maxP := ich.senkouBPeriod
	if ich.kijunPeriod > maxP {
		maxP = ich.kijunPeriod
	}
	ich.highs = make([]float64, maxP)
	ich.lows = make([]float64, maxP)
	ich.index = 0
	ich.count = 0
}

func (ich *Ichimoku) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Ichimoku Cloud",
		Parameters: []indicators.Parameter{
			{Name: "Tenkan", DefaultValue: float64(ich.tenkanPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Kijun", DefaultValue: float64(ich.kijunPeriod), Min: 2, Max: 100, Step: 1},
			{Name: "Senkou B", DefaultValue: float64(ich.senkouBPeriod), Min: 2, Max: 200, Step: 1},
			{Name: "Displacement", DefaultValue: float64(ich.displacement), Min: 1, Max: 100, Step: 1},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "Tenkan-sen", Color: "#2196F3", Style: indicators.StyleLine, Width: 1},
			{Name: "Kijun-sen", Color: "#F44336", Style: indicators.StyleLine, Width: 1},
			{Name: "Senkou A", Color: "#4CAF50", Style: indicators.StyleLine, Width: 1},
			{Name: "Senkou B", Color: "#FF9800", Style: indicators.StyleLine, Width: 1},
			{Name: "Chikou", Color: "#9C27B0", Style: indicators.StyleLine, Width: 1},
		},
	}
}
