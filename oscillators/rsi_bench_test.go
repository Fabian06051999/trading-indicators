package oscillators

import (
	"math/rand"
	"testing"

	"github.com/Fabian06051999/trading-indicators"
)

func generateCandles(n int) []indicators.OHLCV {
	candles := make([]indicators.OHLCV, n)
	price := 100.0
	for i := range candles {
		change := (rand.Float64() - 0.5) * 2
		price += change
		candles[i] = indicators.OHLCV{
			Timestamp: int64(i),
			Open:      price - 0.5,
			High:      price + rand.Float64(),
			Low:       price - rand.Float64(),
			Close:     price,
			Volume:    rand.Float64() * 10000,
		}
	}
	return candles
}

func BenchmarkRSI_1M_Batch(b *testing.B) {
	candles := generateCandles(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rsi := NewRSI(14)
		rsi.CalculateAll(candles)
	}
}

func BenchmarkRSI_1M_Incremental(b *testing.B) {
	candles := generateCandles(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rsi := NewRSI(14)
		for _, c := range candles {
			rsi.UpdateAll(c)
		}
	}
}

func BenchmarkRSI_SingleUpdate(b *testing.B) {
	rsi := NewRSI(14)
	candles := generateCandles(100)
	rsi.CalculateAll(candles)
	candle := candles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rsi.UpdateAll(candle)
	}
}
