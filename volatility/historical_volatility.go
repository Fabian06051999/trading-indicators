package volatility

import (
	"math"

	"github.com/Fabian06051999/trading-indicators"
)

// HistoricalVolatility implements annualized historical volatility.
type HistoricalVolatility struct {
	period    int
	annFactor float64
	returns   []float64
	prevClose float64
	index     int
	count     int
	out       []float64
}

func NewHistoricalVolatility(period int, annualizationFactor float64) *HistoricalVolatility {
	return &HistoricalVolatility{
		period:    period,
		annFactor: annualizationFactor,
		returns:   make([]float64, period),
		out:       make([]float64, 1),
	}
}

func (h *HistoricalVolatility) CalculateAll(candles []indicators.OHLCV) [][]float64 {
	result := make([]float64, len(candles))
	h.Reset()

	for i, c := range candles {
		h.update(c)
		result[i] = h.out[0]
	}
	return [][]float64{result}
}

func (h *HistoricalVolatility) UpdateAll(candle indicators.OHLCV) []float64 {
	h.update(candle)
	return h.out
}

func (h *HistoricalVolatility) update(candle indicators.OHLCV) {
	h.count++

	if h.count == 1 {
		h.prevClose = candle.Close
		h.out[0] = math.NaN()
		return
	}

	logReturn := math.Log(candle.Close / h.prevClose)
	h.prevClose = candle.Close

	h.returns[h.index] = logReturn
	h.index = (h.index + 1) % h.period

	if h.count <= h.period {
		h.out[0] = math.NaN()
		return
	}

	// Calculate standard deviation of log returns
	sum := 0.0
	for i := 0; i < h.period; i++ {
		sum += h.returns[i]
	}
	mean := sum / float64(h.period)

	variance := 0.0
	for i := 0; i < h.period; i++ {
		diff := h.returns[i] - mean
		variance += diff * diff
	}
	variance /= float64(h.period - 1)

	h.out[0] = math.Sqrt(variance) * math.Sqrt(h.annFactor) * 100
}

func (h *HistoricalVolatility) Reset() {
	h.returns = make([]float64, h.period)
	h.prevClose = 0
	h.index = 0
	h.count = 0
}

func (h *HistoricalVolatility) Config() *indicators.IndicatorConfig {
	return &indicators.IndicatorConfig{
		Name: "Historical Volatility",
		Parameters: []indicators.Parameter{
			{Name: "Period", DefaultValue: float64(h.period), Min: 5, Max: 252, Step: 1},
			{Name: "Annualization", DefaultValue: h.annFactor, Min: 1, Max: 365, Step: 1},
		},
		Pane: indicators.PaneSeparate,
		Outputs: []indicators.OutputConfig{
			{Name: "HV", Color: "#D32F2F", Style: indicators.StyleLine, Width: 2},
		},
	}
}
