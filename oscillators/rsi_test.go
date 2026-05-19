package oscillators

import (
	"math"
	"testing"

	"github.com/Fabian06051999/trading-indicators"
)

func TestRSI_WilderReference(t *testing.T) {
	closes := []float64{
		44.34, 44.09, 43.61, 44.33, 44.83, 45.10, 45.42, 45.84, 46.08,
		45.89, 46.03, 45.61, 46.28, 46.28, 46.00, 46.03, 46.41, 46.22, 45.64,
	}

	candles := make([]indicators.OHLCV, len(closes))
	for i, c := range closes {
		candles[i] = indicators.OHLCV{Close: c, High: c + 0.5, Low: c - 0.5, Volume: 1000}
	}

	rsi := NewRSI(14)
	results := rsi.CalculateAll(candles)[0]

	for i := 0; i < 14; i++ {
		if !math.IsNaN(results[i]) {
			t.Errorf("Expected NaN during warmup at index %d, got %f", i, results[i])
		}
	}

	expected := 66.93
	if math.Abs(results[14]-expected) > 0.5 {
		t.Errorf("First RSI value: expected ~%.2f, got %.2f", expected, results[14])
	}

	for i := 14; i < len(results); i++ {
		if results[i] < 0 || results[i] > 100 {
			t.Errorf("RSI out of range at index %d: %f", i, results[i])
		}
	}

	t.Logf("RSI values: %v", results[14:])
}

func TestRSI_AllGains(t *testing.T) {
	candles := make([]indicators.OHLCV, 20)
	for i := range candles {
		price := 100.0 + float64(i)
		candles[i] = indicators.OHLCV{Close: price, High: price, Low: price, Volume: 1000}
	}

	rsi := NewRSI(14)
	results := rsi.CalculateAll(candles)[0]

	lastRSI := results[len(results)-1]
	if lastRSI != 100 {
		t.Errorf("All gains should give RSI=100, got %f", lastRSI)
	}
}

func TestRSI_AllLosses(t *testing.T) {
	candles := make([]indicators.OHLCV, 20)
	for i := range candles {
		price := 200.0 - float64(i)
		candles[i] = indicators.OHLCV{Close: price, High: price, Low: price, Volume: 1000}
	}

	rsi := NewRSI(14)
	results := rsi.CalculateAll(candles)[0]

	lastRSI := results[len(results)-1]
	if lastRSI != 0 {
		t.Errorf("All losses should give RSI=0, got %f", lastRSI)
	}
}

func TestRSI_Incremental(t *testing.T) {
	closes := []float64{
		44.34, 44.09, 43.61, 44.33, 44.83, 45.10, 45.42, 45.84, 46.08,
		45.89, 46.03, 45.61, 46.28, 46.28, 46.00, 46.03, 46.41, 46.22, 45.64,
	}

	candles := make([]indicators.OHLCV, len(closes))
	for i, c := range closes {
		candles[i] = indicators.OHLCV{Close: c, High: c + 0.5, Low: c - 0.5, Volume: 1000}
	}

	rsiBatch := NewRSI(14)
	batchResults := rsiBatch.CalculateAll(candles)[0]

	rsiInc := NewRSI(14)
	for i, c := range candles {
		incResult := rsiInc.UpdateAll(c)[0]
		if math.IsNaN(batchResults[i]) && math.IsNaN(incResult) {
			continue
		}
		if math.Abs(incResult-batchResults[i]) > 0.0001 {
			t.Errorf("Batch vs Incremental mismatch at %d: batch=%f inc=%f", i, batchResults[i], incResult)
		}
	}
}

func TestRSI_Reset(t *testing.T) {
	rsi := NewRSI(14)

	candles := make([]indicators.OHLCV, 20)
	for i := range candles {
		candles[i] = indicators.OHLCV{Close: 100.0 + float64(i), High: 101, Low: 99, Volume: 1000}
	}

	rsi.CalculateAll(candles)
	rsi.Reset()

	result := rsi.UpdateAll(candles[0])[0]
	if !math.IsNaN(result) {
		t.Errorf("After reset, first UpdateAll should return NaN, got %f", result)
	}
}
