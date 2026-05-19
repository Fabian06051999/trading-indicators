package oscillators

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// UltimateOscillator implements Larry Williams' Ultimate Oscillator.
type UltimateOscillator struct {
	period1   int
	period2   int
	period3   int
	bp        []float64
	tr        []float64
	prevClose float64
	index     int
	count     int
	maxPeriod int
}

func NewUltimateOscillator(period1, period2, period3 int) *UltimateOscillator {
	maxP := period3
	if period2 > maxP {
		maxP = period2
	}
	if period1 > maxP {
		maxP = period1
	}
	return &UltimateOscillator{
		period1:   period1,
		period2:   period2,
		period3:   period3,
		bp:        make([]float64, maxP),
		tr:        make([]float64, maxP),
		maxPeriod: maxP,
	}
}

func (u *UltimateOscillator) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	u.Reset()

	for i, c := range candles {
		result[i] = u.UpdateAll(c)[0]
	}
	return [][]float64{result}
}

func (u *UltimateOscillator) UpdateAll(candle indicators.OHLCV) []float64 {
	u.count++

	if u.count == 1 {
		u.prevClose = candle.Close
		return []float64{math.NaN()}
	}

	// Buying Pressure = Close - Min(Low, PrevClose)
	bpVal := candle.Close - math.Min(candle.Low, u.prevClose)
	// True Range = Max(High, PrevClose) - Min(Low, PrevClose)
	trVal := math.Max(candle.High, u.prevClose) - math.Min(candle.Low, u.prevClose)

	u.prevClose = candle.Close

	u.bp[u.index] = bpVal
	u.tr[u.index] = trVal
	u.index = (u.index + 1) % u.maxPeriod

	if u.count <= u.maxPeriod {
		return []float64{math.NaN()}
	}

	// Sum BP and TR for each period
	sumBP1, sumTR1 := u.sumRange(u.period1)
	sumBP2, sumTR2 := u.sumRange(u.period2)
	sumBP3, sumTR3 := u.sumRange(u.period3)

	avg1 := 0.0
	avg2 := 0.0
	avg3 := 0.0
	if sumTR1 != 0 {
		avg1 = sumBP1 / sumTR1
	}
	if sumTR2 != 0 {
		avg2 = sumBP2 / sumTR2
	}
	if sumTR3 != 0 {
		avg3 = sumBP3 / sumTR3
	}

	return []float64{(4*avg1 + 2*avg2 + avg3) / 7.0 * 100}
}

func (u *UltimateOscillator) sumRange(period int) (float64, float64) {
	sumBP := 0.0
	sumTR := 0.0
	for i := 0; i < period; i++ {
		idx := (u.index - 1 - i + u.maxPeriod) % u.maxPeriod
		sumBP += u.bp[idx]
		sumTR += u.tr[idx]
	}
	return sumBP, sumTR
}

func (u *UltimateOscillator) Reset() {
	u.bp = make([]float64, u.maxPeriod)
	u.tr = make([]float64, u.maxPeriod)
	u.prevClose = 0
	u.index = 0
	u.count = 0
}

func (u *UltimateOscillator) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Ultimate Oscillator",
		Parameters: []indicators.Parameter{
			{Name: "Period 1", DefaultValue: float64(u.period1), Min: 2, Max: 50, Step: 1},
			{Name: "Period 2", DefaultValue: float64(u.period2), Min: 2, Max: 100, Step: 1},
			{Name: "Period 3", DefaultValue: float64(u.period3), Min: 2, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "UO",
				Color:  "#1976D2",
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
