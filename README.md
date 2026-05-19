# Trading Indicators (Go Library)

Pure calculation library for technical trading indicators. No server, no frontend — just import and feed it OHLCV data.

## Installation

```bash
go get github.com/Fabian06051999/trading-indicators
```

## Quickstart

```go
package main

import (
    "fmt"
    "github.com/Fabian06051999/trading-indicators"
    "github.com/Fabian06051999/trading-indicators/moving_averages"
    "github.com/Fabian06051999/trading-indicators/oscillators"
)

func main() {
    candles := []indicators.OHLCV{
        {Timestamp: 1, Open: 44, High: 45, Low: 43, Close: 44.5, Volume: 1000},
        {Timestamp: 2, Open: 44.5, High: 46, Low: 44, Close: 45.2, Volume: 1200},
        {Timestamp: 3, Open: 45.2, High: 47, Low: 45, Close: 46.8, Volume: 900},
        // ... more candles
    }

    // --- Batch calculation (all at once) ---
    sma := moving_averages.NewSMA(14)
    values := sma.CalculateAll(candles) // returns [][]float64 (one row per output)
    fmt.Println("SMA values:", values[0])

    // --- Incremental (live updates via WebSocket) ---
    rsi := oscillators.NewRSI(14)
    // First load history
    rsi.CalculateAll(candles)
    // Then update live
    newCandle := indicators.OHLCV{Close: 47.3, High: 48, Low: 46.5, Volume: 800}
    latestRSI := rsi.UpdateAll(newCandle) // returns []float64 (one value per output)
    fmt.Println("Latest RSI:", latestRSI[0])
}
```

## Display Config for SciChart

Every indicator provides a config with defaults for frontend rendering:

```go
rsi := oscillators.NewRSI(14)
cfg := rsi.Config()

fmt.Println(cfg.Name)              // "Relative Strength Index"
fmt.Println(cfg.Pane)              // PaneSeparate (own sub-chart)
fmt.Println(cfg.Outputs[0].Color)  // "#7B1FA2"
fmt.Println(cfg.Outputs[0].YRange) // {Min: 0, Max: 100}
fmt.Println(cfg.Outputs[0].Levels) // [{70, "Overbought", "#EF5350"}, {30, "Oversold", "#66BB6A"}]
```

**Pane types:**
| Value | Meaning | Examples |
|-------|---------|----------|
| `PaneOverlay` | Drawn on the price chart | SMA, EMA, Bollinger, VWAP, Parabolic SAR |
| `PaneSeparate` | Gets its own sub-chart | RSI, MACD, Stochastic, ADX, ATR, OBV |

**Customize colors/style:**
```go
cfg := rsi.Config()
cfg.Outputs[0].Color = "#00FF00"               // change color
cfg.Outputs[0].Width = 3                       // thicker line
cfg.Outputs[0].Style = indicators.StyleDashed  // dashed
cfg.Outputs[0].Levels[0].Value = 80            // adjust overbought level
```

## Unified Interface

Every indicator implements the same interface — no type assertions needed:

```go
type Indicator interface {
    CalculateAll(candles []OHLCV) [][]float64  // batch mode
    UpdateAll(candle OHLCV) []float64           // incremental mode
    Reset()
    Config() *IndicatorConfig
}
```

Single-output indicators (SMA, RSI, ATR) return one row, multi-output (MACD, Bollinger, Ichimoku) return multiple:

```go
// Single output — access [0]
sma := moving_averages.NewSMA(14)
values := sma.CalculateAll(candles)[0]

// Multi output — each index is a series
macd := trend.NewMACD(12, 26, 9)
results := macd.CalculateAll(candles)
// results[0] = MACD Line
// results[1] = Signal Line
// results[2] = Histogram
```

## Available Indicators (~90 total)

### Moving Averages (Overlay) — `moving_averages/`
| Indicator | Constructor | Parameters |
|-----------|-------------|------------|
| SMA | `NewSMA(period)` | Period |
| EMA | `NewEMA(period)` | Period |
| WMA | `NewWMA(period)` | Period |
| DEMA | `NewDEMA(period)` | Period |
| TEMA | `NewTEMA(period)` | Period |
| HMA | `NewHMA(period)` | Period |
| VWMA | `NewVWMA(period)` | Period |
| SMMA | `NewSMMA(period)` | Period |
| KAMA | `NewKAMA(period, fast, slow)` | Period, Fast, Slow |
| ZLEMA | `NewZLEMA(period)` | Period |
| ALMA | `NewALMA(period, offset, sigma)` | Period, Offset, Sigma |
| LSMA | `NewLSMA(period)` | Period |
| McGinley Dynamic | `NewMcGinleyDynamic(period)` | Period |
| Tillson T3 | `NewT3(period, volumeFactor)` | Period, Volume Factor |
| VIDYA | `NewVIDYA(period, cmoPeriod)` | Period, CMO Period |
| FRAMA | `NewFRAMA(period)` | Period |

### Oscillators (Separate Pane) — `oscillators/`
| Indicator | Constructor | Parameters |
|-----------|-------------|------------|
| RSI | `NewRSI(period)` | Period |
| Stochastic | `NewStochastic(k, d, slowing)` | K, D, Slowing |
| Stochastic RSI | `NewStochRSI(rsi, k, d)` | RSI Period, K, D |
| CCI | `NewCCI(period)` | Period |
| Williams %R | `NewWilliamsR(period)` | Period |
| MFI | `NewMFI(period)` | Period |
| CMO | `NewCMO(period)` | Period |
| DeMarker | `NewDeMarker(period)` | Period |
| Awesome Oscillator | `NewAwesomeOscillator()` | — |
| Accelerator Osc. | `NewAcceleratorOscillator()` | — |
| Ultimate Oscillator | `NewUltimateOscillator(p1, p2, p3)` | Period 1/2/3 |
| PPO | `NewPPO(fast, slow, signal)` | Fast, Slow, Signal |

### Trend — `trend/`
| Indicator | Constructor | Pane |
|-----------|-------------|------|
| MACD | `NewMACD(fast, slow, signal)` | Separate |
| ADX (+DI/-DI) | `NewADX(period)` | Separate |
| Parabolic SAR | `NewParabolicSAR(step, max)` | Overlay |
| Supertrend | `NewSupertrend(period, multiplier)` | Overlay |
| Aroon (Up/Down) | `NewAroon(period)` | Separate |
| Vortex (VI+/VI-) | `NewVortex(period)` | Separate |
| Ichimoku Cloud | `NewIchimoku(tenkan, kijun, senkouB, disp)` | Overlay |
| TRIX | `NewTRIX(period)` | Separate |
| DPO | `NewDPO(period)` | Separate |
| Mass Index | `NewMassIndex(emaPeriod, sumPeriod)` | Separate |
| MA Envelope | `NewEMAEnvelope(period, percentage)` | Overlay |

### Volatility — `volatility/`
| Indicator | Constructor | Pane |
|-----------|-------------|------|
| Bollinger Bands | `NewBollingerBands(period, stdDev)` | Overlay |
| ATR | `NewATR(period)` | Separate |
| ATR Percent | `NewATRPercent(period)` | Separate |
| Keltner Channel | `NewKeltnerChannel(ema, atr, multi)` | Overlay |
| Donchian Channel | `NewDonchianChannel(period)` | Overlay |
| Standard Deviation | `NewStdDev(period)` | Separate |
| Historical Volatility | `NewHistoricalVolatility(period, annFactor)` | Separate |
| Chaikin Volatility | `NewChaikinVolatility(ema, roc)` | Separate |
| Ulcer Index | `NewUlcerIndex(period)` | Separate |

### Volume — `volume/`
| Indicator | Constructor | Pane |
|-----------|-------------|------|
| OBV | `NewOBV()` | Separate |
| VWAP | `NewVWAP()` | Overlay |
| A/D Line | `NewADLine()` | Separate |
| Chaikin Money Flow | `NewChaikinMF(period)` | Separate |
| Force Index | `NewForceIndex(period)` | Separate |
| Volume Oscillator | `NewVolumeOscillator(fast, slow)` | Separate |
| Ease of Movement | `NewEaseOfMovement(period)` | Separate |
| Negative Volume Index | `NewNVI()` | Separate |
| Positive Volume Index | `NewPVI()` | Separate |
| Volume Profile (POC) | `NewVolumeProfile(period, bins)` | Overlay |

### Momentum — `momentum/`
| Indicator | Constructor | Pane |
|-----------|-------------|------|
| ROC | `NewROC(period)` | Separate |
| Momentum | `NewMomentum(period)` | Separate |
| TSI | `NewTSI(long, short, signal)` | Separate |
| KST | `NewKST(roc1..4, sma1..4, signal)` | Separate |
| Coppock Curve | `NewCoppockCurve(wma, roc1, roc2)` | Separate |
| Elder Ray | `NewElderRay(period)` | Separate |
| Williams AD | `NewWilliamsAD()` | Separate |
| Chande Forecast | `NewChandeForecast(period)` | Separate |
| Bollinger %B | `NewPercentB(period, stdDev)` | Separate |
| Pivot Points | `NewPivotPoints()` | Overlay |
| RVI | `NewRVI(period)` | Separate |
| Connors RSI | `NewConnorsRSI(rsi, streak, rank)` | Separate |
| Schaff Trend Cycle | `NewSchaffTrendCycle(fast, slow, cycle)` | Separate |
| Squeeze Momentum | `NewSqueezeMomentum(bbP, bbSD, kcP, kcM)` | Separate |
| Fisher Transform | `NewFisherTransform(period)` | Separate |
| Choppiness Index | `NewChoppinessIndex(period)` | Separate |
| WaveTrend | `NewWaveTrend(chPeriod, avgPeriod)` | Separate |

## Thread Safety

Indicators are **NOT** thread-safe. Each goroutine must use its own instance:

```go
// WRONG — race condition
rsi := oscillators.NewRSI(14)
go func() { rsi.UpdateAll(candle1) }()
go func() { rsi.UpdateAll(candle2) }()

// CORRECT — one instance per goroutine
go func() {
    rsi := oscillators.NewRSI(14)
    rsi.UpdateAll(candle1)
}()
```

This is by design: no mutex overhead on the hot path. If you process multiple symbols, create one indicator instance per symbol.

## Design Decisions

**Why Batch + Incremental?**
Historical data arrives as a block from the server → `Calculate()`/`CalculateAll()`. Live data arrives one candle at a time via WebSocket → `Update()`/`UpdateAll()`. You can combine both: batch for history load, then incremental for live.

**Why is Config separate from calculation?**
Calculation is pure math. Config (colors, style, levels) is just metadata for SciChart. Changing a color never triggers recalculation.

**Why Ring Buffers instead of arrays?**
Performance at 1M+ candles. A ring buffer uses O(1) memory regardless of how many candles pass through. No copying, no growing.

## Benchmarks (AMD Ryzen 7, 1M candles)

| Indicator | 1M Batch | Single Update | Allocs/Update |
|-----------|----------|---------------|---------------|
| SMA(200) | 9.6ms | 3.8ns | 0 |
| EMA(200) | 6.2ms | ~5ns | 0 |
| RSI(14) | 41ms | 45ns | 1 |

Run benchmarks: `go test ./... -bench=. -benchmem`

**Why `NaN` for "not ready yet"?**
Indicators need a warmup phase (e.g., SMA(14) needs 14 candles). Until then they return `math.NaN()`. In the frontend: `if math.IsNaN(value) { skip }`. This avoids ambiguity — a real `0` (e.g., ROC returning zero) is never confused with "not enough data".
