package trend

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// ADX implements the Average Directional Index with +DI and -DI.
type ADX struct {
	period    int
	prevHigh  float64
	prevLow   float64
	prevClose float64
	smoothPDI float64
	smoothNDI float64
	smoothTR  float64
	adxSum    float64
	adxValue  float64
	count     int
}

func NewADX(period int) *ADX {
	return &ADX{
		period: period,
	}
}

func (a *ADX) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	adxValues := make([]float64, len(candles))
	pdiValues := make([]float64, len(candles))
	ndiValues := make([]float64, len(candles))
	a.Reset()

	for i, c := range candles {
		values := a.UpdateAll(c)
		adxValues[i] = values[0]
		pdiValues[i] = values[1]
		ndiValues[i] = values[2]
	}
	return [][]float64{adxValues, pdiValues, ndiValues}
}

func (a *ADX) UpdateAll(candle indicators.OHLCV) []float64 {
	a.count++

	if a.count == 1 {
		a.prevHigh = candle.High
		a.prevLow = candle.Low
		a.prevClose = candle.Close
		return []float64{0, 0, 0}
	}

	// True Range
	tr := math.Max(candle.High-candle.Low,
		math.Max(math.Abs(candle.High-a.prevClose), math.Abs(candle.Low-a.prevClose)))

	// Directional Movement
	upMove := candle.High - a.prevHigh
	downMove := a.prevLow - candle.Low

	pDM := 0.0
	nDM := 0.0
	if upMove > downMove && upMove > 0 {
		pDM = upMove
	}
	if downMove > upMove && downMove > 0 {
		nDM = downMove
	}

	a.prevHigh = candle.High
	a.prevLow = candle.Low
	a.prevClose = candle.Close

	p := float64(a.period)

	if a.count <= a.period+1 {
		a.smoothTR += tr
		a.smoothPDI += pDM
		a.smoothNDI += nDM

		if a.count == a.period+1 {
			// First smoothed values
			pdi := 0.0
			ndi := 0.0
			if a.smoothTR != 0 {
				pdi = (a.smoothPDI / a.smoothTR) * 100
				ndi = (a.smoothNDI / a.smoothTR) * 100
			}
			dx := 0.0
			if pdi+ndi != 0 {
				dx = math.Abs(pdi-ndi) / (pdi + ndi) * 100
			}
			a.adxSum = dx
			return []float64{0, pdi, ndi}
		}
		return []float64{0, 0, 0}
	}

	// Smoothed values using Wilder's method
	a.smoothTR = a.smoothTR - a.smoothTR/p + tr
	a.smoothPDI = a.smoothPDI - a.smoothPDI/p + pDM
	a.smoothNDI = a.smoothNDI - a.smoothNDI/p + nDM

	pdi := 0.0
	ndi := 0.0
	if a.smoothTR != 0 {
		pdi = (a.smoothPDI / a.smoothTR) * 100
		ndi = (a.smoothNDI / a.smoothTR) * 100
	}

	dx := 0.0
	if pdi+ndi != 0 {
		dx = math.Abs(pdi-ndi) / (pdi + ndi) * 100
	}

	step := a.count - a.period - 1
	if step < a.period {
		a.adxSum += dx
		if step == a.period-1 {
			a.adxValue = a.adxSum / p
		}
		return []float64{0, pdi, ndi}
	}

	a.adxValue = (a.adxValue*float64(a.period-1) + dx) / p
	return []float64{a.adxValue, pdi, ndi}
}

func (a *ADX) Reset() {
	a.prevHigh = 0
	a.prevLow = 0
	a.prevClose = 0
	a.smoothPDI = 0
	a.smoothNDI = 0
	a.smoothTR = 0
	a.adxSum = 0
	a.adxValue = 0
	a.count = 0
}

func (a *ADX) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Average Directional Index",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(a.period), Min: 2, Max: 100, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "ADX",
				Color:  "#FF5722",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 25, Label: "Trend", Color: "#9E9E9E"},
				},
			},
			{Name: "+DI", Color: "#4CAF50", Style: indicators.StyleLine, Width: 1},
			{Name: "-DI", Color: "#F44336", Style: indicators.StyleLine, Width: 1},
		},
	}
}
