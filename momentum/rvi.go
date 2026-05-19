package momentum

import (
	"math"
	"github.com/Fabian06051999/trading-indicators"
)

// RVI implements the Relative Vigor Index.
type RVI struct {
	period  int
	numBuf  []float64
	denBuf  []float64
	sigBuf  []float64
	opens   []float64
	highs   []float64
	lows    []float64
	closes  []float64
	index   int
	count   int
}

func NewRVI(period int) *RVI {
	return &RVI{
		period: period,
		numBuf: make([]float64, period),
		denBuf: make([]float64, period),
		sigBuf: make([]float64, 4),
		opens:  make([]float64, 4),
		highs:  make([]float64, 4),
		lows:   make([]float64, 4),
		closes: make([]float64, 4),
	}
}

func (r *RVI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	rviValues := make([]float64, len(candles))
	sigValues := make([]float64, len(candles))
	r.Reset()

	for i, c := range candles {
		values := r.UpdateAll(c)
		rviValues[i] = values[0]
		sigValues[i] = values[1]
	}
	return [][]float64{rviValues, sigValues}
}

func (r *RVI) UpdateAll(candle indicators.OHLCV) []float64 {
	r.count++

	// Shift OHLC buffers
	copy(r.opens[0:3], r.opens[1:4])
	copy(r.highs[0:3], r.highs[1:4])
	copy(r.lows[0:3], r.lows[1:4])
	copy(r.closes[0:3], r.closes[1:4])
	r.opens[3] = candle.Open
	r.highs[3] = candle.High
	r.lows[3] = candle.Low
	r.closes[3] = candle.Close

	if r.count < 4 {
		return []float64{math.NaN(), math.NaN()}
	}

	// Symmetrically weighted moving average of numerator and denominator
	num := (r.closes[3]-r.opens[3] + 2*(r.closes[2]-r.opens[2]) + 2*(r.closes[1]-r.opens[1]) + (r.closes[0]-r.opens[0])) / 6.0
	den := (r.highs[3]-r.lows[3] + 2*(r.highs[2]-r.lows[2]) + 2*(r.highs[1]-r.lows[1]) + (r.highs[0]-r.lows[0])) / 6.0

	idx := (r.count - 4) % r.period
	r.numBuf[idx] = num
	r.denBuf[idx] = den

	if r.count < r.period+3 {
		return []float64{math.NaN(), math.NaN()}
	}

	sumNum := 0.0
	sumDen := 0.0
	for i := 0; i < r.period; i++ {
		sumNum += r.numBuf[i]
		sumDen += r.denBuf[i]
	}

	rviVal := 0.0
	if sumDen != 0 {
		rviVal = sumNum / sumDen
	}

	// Signal line (symmetrically weighted MA of RVI)
	copy(r.sigBuf[0:3], r.sigBuf[1:4])
	r.sigBuf[3] = rviVal

	signal := (r.sigBuf[3] + 2*r.sigBuf[2] + 2*r.sigBuf[1] + r.sigBuf[0]) / 6.0

	return []float64{rviVal, signal}
}

func (r *RVI) Reset() {
	r.numBuf = make([]float64, r.period)
	r.denBuf = make([]float64, r.period)
	r.sigBuf = make([]float64, 4)
	r.opens = make([]float64, 4)
	r.highs = make([]float64, 4)
	r.lows = make([]float64, 4)
	r.closes = make([]float64, 4)
	r.index = 0
	r.count = 0
}

func (r *RVI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Relative Vigor Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(r.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "RVI", Color: "#1565C0", Style: indicators.StyleLine, Width: 2},
			{Name: "Signal", Color: "#FF6F00", Style: indicators.StyleDashed, Width: 1},
		},
	}
}
