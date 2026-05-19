package moving_averages

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

func BenchmarkSMA_1M_Batch(b *testing.B) {
	candles := generateCandles(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sma := NewSMA(200)
		sma.CalculateAll(candles)
	}
}

func BenchmarkEMA_1M_Batch(b *testing.B) {
	candles := generateCandles(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ema := NewEMA(200)
		ema.CalculateAll(candles)
	}
}

func BenchmarkHMA_1M_Batch(b *testing.B) {
	candles := generateCandles(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hma := NewHMA(200)
		hma.CalculateAll(candles)
	}
}

func BenchmarkSMA_SingleUpdate(b *testing.B) {
	sma := NewSMA(200)
	candles := generateCandles(300)
	sma.CalculateAll(candles)
	candle := candles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sma.UpdateAll(candle)
	}
}
