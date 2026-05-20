package trend

import (
	"github.com/Fabian06051999/trading-indicators"
)

// ParabolicSAR implements the Parabolic Stop and Reverse.
type ParabolicSAR struct {
	afStep   float64
	afMax    float64
	af       float64
	ep       float64
	sar      float64
	isLong   bool
	count    int
	prevHigh float64
	prevLow  float64
	out      []float64
}

func NewParabolicSAR(afStep, afMax float64) *ParabolicSAR {
	return &ParabolicSAR{
		afStep: afStep,
		afMax:  afMax,
		out:    make([]float64, 1),
	}
}

func (p *ParabolicSAR) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	p.Reset()

	for i, c := range candles {
		p.update(c)
		result[i] = p.out[0]
	}
	return [][]float64{result}
}

func (p *ParabolicSAR) UpdateAll(candle indicators.OHLCV) []float64 {
	p.update(candle)
	return p.out
}

func (p *ParabolicSAR) update(candle indicators.OHLCV) {
	p.count++

	if p.count == 1 {
		p.prevHigh = candle.High
		p.prevLow = candle.Low
		p.sar = candle.Low
		p.ep = candle.High
		p.af = p.afStep
		p.isLong = true
		p.out[0] = p.sar
		return
	}

	if p.count == 2 {
		if candle.Close > p.prevHigh {
			p.isLong = true
			p.sar = p.prevLow
			p.ep = candle.High
		} else {
			p.isLong = false
			p.sar = p.prevHigh
			p.ep = candle.Low
		}
		p.af = p.afStep
		p.prevHigh = candle.High
		p.prevLow = candle.Low
		p.out[0] = p.sar
		return
	}

	newSAR := p.sar + p.af*(p.ep-p.sar)

	if p.isLong {
		if newSAR > candle.Low {
			// Reverse to short
			p.isLong = false
			newSAR = p.ep
			p.ep = candle.Low
			p.af = p.afStep
		} else {
			if candle.High > p.ep {
				p.ep = candle.High
				p.af += p.afStep
				if p.af > p.afMax {
					p.af = p.afMax
				}
			}
			// SAR must not be above prior two lows
			if newSAR > p.prevLow {
				newSAR = p.prevLow
			}
		}
	} else {
		if newSAR < candle.High {
			// Reverse to long
			p.isLong = true
			newSAR = p.ep
			p.ep = candle.High
			p.af = p.afStep
		} else {
			if candle.Low < p.ep {
				p.ep = candle.Low
				p.af += p.afStep
				if p.af > p.afMax {
					p.af = p.afMax
				}
			}
			// SAR must not be below prior two highs
			if newSAR < p.prevHigh {
				newSAR = p.prevHigh
			}
		}
	}

	p.sar = newSAR
	p.prevHigh = candle.High
	p.prevLow = candle.Low
	p.out[0] = p.sar
}

func (p *ParabolicSAR) Reset() {
	p.af = p.afStep
	p.ep = 0
	p.sar = 0
	p.isLong = true
	p.count = 0
	p.prevHigh = 0
	p.prevLow = 0
}

func (p *ParabolicSAR) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Parabolic SAR",
		Parameters: []indicators.Parameter{
			{Name: "AF Step", DefaultValue: p.afStep, Min: 0.01, Max: 0.1, Step: 0.01},
			{Name: "AF Max", DefaultValue: p.afMax, Min: 0.1, Max: 0.5, Step: 0.01},
		},
		Pane: indicators.PaneOverlay,
		Outputs: []indicators.OutputConfig{
			{Name: "SAR", Color: "#FF9800", Style: indicators.StyleDots, Width: 2},
		},
	}
}
