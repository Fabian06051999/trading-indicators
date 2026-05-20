package momentum

import (
	"github.com/Fabian06051999/trading-indicators"
	"github.com/Fabian06051999/trading-indicators/oscillators"
	"math"
)

// ConnorsRSI implements the Connors RSI (combination of RSI, streak RSI, and percentile rank).
type ConnorsRSI struct {
	out          []float64
	rsiPeriod    int
	streakPeriod int
	rankPeriod   int
	rsi          *oscillators.RSI
	streakRSI    *oscillators.RSI
	streak       int
	prevClose    float64
	rocBuffer    []float64
	rocIndex     int
	count        int
}

func NewConnorsRSI(rsiPeriod, streakPeriod, rankPeriod int) *ConnorsRSI {
	return &ConnorsRSI{
		rsiPeriod:    rsiPeriod,
		streakPeriod: streakPeriod,
		rankPeriod:   rankPeriod,
		rsi:          oscillators.NewRSI(rsiPeriod),
		streakRSI:    oscillators.NewRSI(streakPeriod),
		rocBuffer:    make([]float64, rankPeriod),
		out:          make([]float64, 1),
	}
}

func (c *ConnorsRSI) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	c.Reset()

	for i, candle := range candles {
		result[i] = c.UpdateAll(candle)[0]
	}
	return [][]float64{result}
}

func (c *ConnorsRSI) UpdateAll(candle indicators.OHLCV) []float64 {
	c.count++

	// Standard RSI
	rsiVal := c.rsi.UpdateAll(candle)[0]

	// Streak calculation
	if c.count > 1 {
		if candle.Close > c.prevClose {
			if c.streak >= 0 {
				c.streak++
			} else {
				c.streak = 1
			}
		} else if candle.Close < c.prevClose {
			if c.streak <= 0 {
				c.streak--
			} else {
				c.streak = -1
			}
		} else {
			c.streak = 0
		}
	}

	// Streak RSI
	streakRSIVal := c.streakRSI.UpdateAll(indicators.OHLCV{Close: float64(c.streak)})[0]

	// Percentile rank of ROC
	roc := 0.0
	if c.prevClose != 0 {
		roc = ((candle.Close - c.prevClose) / c.prevClose) * 100
	}
	c.prevClose = candle.Close

	c.rocBuffer[c.rocIndex] = roc
	c.rocIndex = (c.rocIndex + 1) % c.rankPeriod

	if c.count < c.rankPeriod {
		c.out[0] = math.NaN()
		return c.out
	}

	// Count how many past ROC values are below current
	below := 0
	for i := 0; i < c.rankPeriod; i++ {
		if c.rocBuffer[i] < roc {
			below++
		}
	}
	percentRank := (float64(below) / float64(c.rankPeriod)) * 100

	if rsiVal == 0 || streakRSIVal == 0 {
		c.out[0] = math.NaN()
		return c.out
	}

	c.out[0] = (rsiVal + streakRSIVal + percentRank) / 3.0
	return c.out
}

func (c *ConnorsRSI) Reset() {
	c.rsi.Reset()
	c.streakRSI.Reset()
	c.streak = 0
	c.prevClose = 0
	c.rocBuffer = make([]float64, c.rankPeriod)
	c.rocIndex = 0
	c.count = 0
}

func (c *ConnorsRSI) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Connors RSI",
		Parameters: []indicators.Parameter{
			{Name: "RSI Period", DefaultValue: float64(c.rsiPeriod), Min: 2, Max: 50, Step: 1},
			{Name: "Streak Period", DefaultValue: float64(c.streakPeriod), Min: 2, Max: 20, Step: 1},
			{Name: "Rank Period", DefaultValue: float64(c.rankPeriod), Min: 10, Max: 200, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{
				Name:   "CRSI",
				Color:  "#880E4F",
				Style:  indicators.StyleLine,
				Width:  2,
				YRange: &indicators.YRange{Min: 0, Max: 100},
				Levels: []indicators.Level{
					{Value: 70, Label: "Overbought", Color: "#EF5350"},
					{Value: 30, Label: "Oversold", Color: "#66BB6A"},
				},
			},
		},
	}
}
