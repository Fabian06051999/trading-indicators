# Trading Indicators

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Indicators](https://img.shields.io/badge/Indicators-75-blue)]()
[![Zero Alloc](https://img.shields.io/badge/Updates-Zero%20Alloc-green)]()

Pure Go calculation library for technical trading indicators. Zero allocations on the hot path. Built for high-frequency crypto data pipelines with SciChart frontend integration.

## Features

- **75 indicators** across 6 categories (Moving Averages, Oscillators, Trend, Volatility, Volume, Momentum)
- **Zero-allocation updates** — pre-allocated output slices, no GC pressure
- **Unified interface** — all indicators implement the same `Indicator` interface
- **Registry** — instantiate any indicator by name with dynamic parameters
- **SciChart-compatible config** — panes, colors, opacity, fill, dash patterns, conditional colors, point markers
- **Batch + Incremental** — full history calculation or single-candle live updates
- **14.5ns per update** (RSI benchmark, AMD Ryzen 7)

## Installation

```bash
go get github.com/Fabian06051999/trading-indicators
```

## Quick Start

```go
package main

import (
    "fmt"
    "math"

    "github.com/Fabian06051999/trading-indicators"
    _ "github.com/Fabian06051999/trading-indicators/oscillators"
    _ "github.com/Fabian06051999/trading-indicators/trend"
)

func main() {
    candles := []indicators.OHLCV{
        {Timestamp: 1, Open: 44, High: 45, Low: 43, Close: 44.5, Volume: 1000},
        {Timestamp: 2, Open: 44.5, High: 46, Low: 44, Close: 45.2, Volume: 1200},
        // ... more candles
    }

    // Create indicator by name (via registry)
    rsi, _ := indicators.Create("RSI", map[string]float64{"Period": 14})

    // Batch: calculate full history
    results := rsi.CalculateAll(candles)[0]
    for _, v := range results {
        if !math.IsNaN(v) {
            fmt.Printf("RSI: %.2f\n", v)
        }
    }

    // Incremental: live updates (zero allocations)
    newCandle := indicators.OHLCV{Close: 47.3, High: 48, Low: 46.5, Volume: 800}
    value := rsi.UpdateAll(newCandle)[0]
    fmt.Printf("Live RSI: %.2f\n", value)
}
```

## Registry

All indicators auto-register via `init()`. Import packages with blank identifier to activate:

```go
import (
    "github.com/Fabian06051999/trading-indicators"
    _ "github.com/Fabian06051999/trading-indicators/moving_averages"
    _ "github.com/Fabian06051999/trading-indicators/oscillators"
    _ "github.com/Fabian06051999/trading-indicators/trend"
    _ "github.com/Fabian06051999/trading-indicators/volatility"
    _ "github.com/Fabian06051999/trading-indicators/volume"
    _ "github.com/Fabian06051999/trading-indicators/momentum"
)

// Create by name — no switch/case needed
macd, _ := indicators.Create("MACD", map[string]float64{"Fast": 12, "Slow": 26, "Signal": 9})

// List all available indicators (e.g. for UI dropdown)
all := indicators.ListAll()

// Filter by category
oscillators := indicators.ListCategory("Oscillators")

// Use defaults (nil params = all defaults)
rsi, _ := indicators.Create("RSI", nil)  // Period defaults to 14
```

## Interface

Every indicator implements one interface:

```go
type Indicator interface {
    CalculateAll(candles []OHLCV) [][]float64  // batch: full history
    UpdateAll(candle OHLCV) []float64           // incremental: one candle
    Reset()                                     // clear state
    Config() *IndicatorConfig                   // display metadata
}
```

Single-output indicators return `[][]float64` with one row. Multi-output (MACD, Bollinger, Ichimoku) return multiple rows:

```go
macd, _ := indicators.Create("MACD", nil)
results := macd.CalculateAll(candles)
// results[0] = MACD Line
// results[1] = Signal Line
// results[2] = Histogram
```

## SciChart Integration

Every indicator carries display metadata that maps directly to SciChart's rendering API:

```go
cfg := rsi.Config()

cfg.Name                    // "Relative Strength Index"
cfg.Pane                    // PaneSeparate (own sub-chart)
cfg.Outputs[0].Color        // "#7B1FA2"
cfg.Outputs[0].Opacity      // 1.0
cfg.Outputs[0].Style        // StyleLine
cfg.Outputs[0].Width        // 2
cfg.Outputs[0].DashArray    // []int{5, 5} for dashed lines
cfg.Outputs[0].FillColor    // "#EF535020" (band fill with alpha)
cfg.Outputs[0].UpColor      // "#4CAF50" (green when rising)
cfg.Outputs[0].DownColor    // "#F44336" (red when falling)
cfg.Outputs[0].Marker       // MarkerEllipse (for dot series)
cfg.Outputs[0].YRange       // &YRange{Min: 0, Max: 100}
cfg.Outputs[0].Levels       // [{Value: 70, Label: "Overbought", Color: "#EF5350"}]
```

| SciChart Concept | Our Field | Example |
|------------------|-----------|---------|
| Sub-chart / Pane | `Pane` | `PaneOverlay` or `PaneSeparate` |
| `stroke` | `Color` | `"#2196F3"` |
| `strokeThickness` | `Width` | 2 |
| `opacity` | `Opacity` | 0.8 |
| `strokeDashArray` | `DashArray` | `[5, 5]` |
| Band fill | `FillColor` | `"#EF535020"` |
| Histogram +/- colors | `UpColor` / `DownColor` | green / red |
| `pointMarker` | `Marker` | `MarkerEllipse` |
| `visibleRange` | `YRange` | `{Min: 0, Max: 100}` |
| Horizontal annotations | `Levels` | Overbought/Oversold lines |

## Available Indicators

### Moving Averages — `moving_averages/`
SMA, EMA, WMA, DEMA, TEMA, HMA, VWMA, SMMA, KAMA, ZLEMA, ALMA, LSMA, McGinley Dynamic, Tillson T3, VIDYA, FRAMA

### Oscillators — `oscillators/`
RSI, Stochastic, Stochastic RSI, CCI, Williams %R, MFI, CMO, DeMarker, Awesome Oscillator, Accelerator Oscillator, Ultimate Oscillator, PPO

### Trend — `trend/`
MACD, ADX (+DI/-DI), Parabolic SAR, Supertrend, Aroon, Vortex, Ichimoku Cloud, TRIX, DPO, Mass Index, MA Envelope

### Volatility — `volatility/`
Bollinger Bands, ATR, ATR Percent, Keltner Channel, Donchian Channel, Standard Deviation, Historical Volatility, Chaikin Volatility, Ulcer Index

### Volume — `volume/`
OBV, VWAP, A/D Line, Chaikin Money Flow, Force Index, Volume Oscillator, Ease of Movement, NVI, PVI, Volume Profile (POC)

### Momentum — `momentum/`
ROC, Momentum, TSI, KST, Coppock Curve, Elder Ray, Williams AD, Chande Forecast, Bollinger %B, Pivot Points, RVI, Connors RSI, Schaff Trend Cycle, Squeeze Momentum, Fisher Transform, Choppiness Index, WaveTrend

## Performance

| Indicator | 1M Candles (Batch) | Single Update | Allocs/Update |
|-----------|-------------------|---------------|---------------|
| SMA(200) | 9.6ms | 3.8ns | 0 |
| EMA(200) | 6.2ms | ~5ns | 0 |
| RSI(14) | 13.9ms | 14.5ns | 0 |

All indicators use **zero-allocation updates**. The returned `[]float64` slice is reused between calls:

```go
value := indicator.UpdateAll(candle)[0]  // read immediately — safe
// Do NOT store the slice reference — it's overwritten on next call
```

Run benchmarks: `go test ./... -bench=. -benchmem`

## Thread Safety

Indicators are **not** thread-safe. Create one instance per goroutine/symbol:

```go
// One indicator per symbol — no locks needed
btcRSI := indicators.MustCreate("RSI", nil)
ethRSI := indicators.MustCreate("RSI", nil)
```

This is by design: zero mutex overhead on the hot path.

## Warmup

Indicators return `math.NaN()` until they have enough data (e.g., SMA(14) needs 14 candles). Check with:

```go
if math.IsNaN(value) {
    continue // skip — not enough data yet
}
```

## Design Decisions

See [ARCHITECTURE.md](ARCHITECTURE.md) for the full design document, including:
- Why unified interface over separate single/multi-output types
- Why zero-alloc with reused slices
- Why NaN instead of 0 for warmup
- Why ring buffers for O(1) memory
- Step-by-step template for adding new indicators

## License

MIT
