# Architecture

This document is for AI assistants (Claude, Cursor, Copilot) and developers who need to understand or extend this library.

## What This Is

A pure Go calculation library providing ~90 technical trading indicators for a crypto trading platform. This library is consumed by a Go backend that feeds calculated values to a SciChart frontend via WebSocket.

**This library does NOT:**
- Serve HTTP/WebSocket endpoints
- Render anything (SciChart handles rendering)
- Manage data pipelines (the consumer handles that)
- Store state between sessions

## Core Interface

Every indicator implements one interface:

```go
type Indicator interface {
    CalculateAll(candles []OHLCV) [][]float64  // batch: full history
    UpdateAll(candle OHLCV) []float64           // incremental: one new candle
    Reset()
    Config() *IndicatorConfig
}
```

**Why unified:** The consumer can treat all indicators identically — no type assertions, no special cases. Single-output indicators return `[][]float64` with one row, multi-output with multiple.

**Why `[]float64` not `float64` for UpdateAll:** Unified interface. MACD returns 3 values, RSI returns 1. Both use the same method signature.

## Zero-Allocation Design

Every indicator pre-allocates its output slice in the constructor:

```go
type RSI struct {
    out []float64  // allocated once: make([]float64, 1)
    // ...
}

func (r *RSI) UpdateAll(candle OHLCV) []float64 {
    r.update(candle)    // writes into r.out[0]
    return r.out        // returns the same slice every time
}
```

**IMPORTANT:** The returned slice is reused. If the consumer needs to store the value, they must copy it:
```go
value := indicator.UpdateAll(candle)[0]  // OK: read immediately
// NOT OK: storing the slice reference (will be overwritten next call)
```

This gives zero allocations per update — critical for crypto tick data (thousands of updates/second).

## Internal Pattern

Each indicator has:
- `UpdateAll()` — public, returns `self.out` (the pre-allocated slice)
- `update()` — private, contains the actual math, writes into `self.out[N]`
- `CalculateAll()` — calls `update()` in a loop, reads from `self.out[0]`

## Warmup & NaN

Indicators need N candles before producing valid output. During warmup, they write `math.NaN()` into their output. The consumer should skip NaN values:

```go
if math.IsNaN(value) { continue }
```

**Why NaN not 0:** Indicators like ROC, Momentum, CMO can legitimately compute 0. NaN is unambiguous.

## Display Config

Each indicator carries metadata for SciChart rendering:

```go
type IndicatorConfig struct {
    Name       string
    Parameters []Parameter     // adjustable inputs (period, multiplier, etc.)
    Pane       PaneType        // PaneOverlay (on price chart) or PaneSeparate (own sub-chart)
    Outputs    []OutputConfig  // one per series (color, style, width, Y-range, levels)
}
```

Config is decoupled from calculation — changing colors never triggers recalculation.

## How to Add a New Indicator

1. Choose the correct package (`oscillators/`, `trend/`, `volatility/`, `volume/`, `momentum/`, `moving_averages/`)
2. Create a new `.go` file
3. Follow this template:

```go
package oscillators

import (
    "math"
    "github.com/Fabian06051999/trading-indicators"
)

type MyIndicator struct {
    period int
    out    []float64
    // ... internal state
}

func NewMyIndicator(period int) *MyIndicator {
    period = indicators.ClampMin(period, 2)
    return &MyIndicator{
        period: period,
        out:    make([]float64, 1),  // number of output series
    }
}

func (m *MyIndicator) CalculateAll(candles []indicators.OHLCV) [][]float64 {
    result := make([]float64, len(candles))
    m.Reset()
    for i, c := range candles {
        m.update(c)
        result[i] = m.out[0]
    }
    return [][]float64{result}
}

func (m *MyIndicator) UpdateAll(candle indicators.OHLCV) []float64 {
    m.update(candle)
    return m.out
}

func (m *MyIndicator) update(candle indicators.OHLCV) {
    // Your math here
    // During warmup: m.out[0] = math.NaN(); return
    // When ready:    m.out[0] = calculatedValue
}

func (m *MyIndicator) Reset() {
    // Reset all internal state (NOT the out slice)
}

func (m *MyIndicator) Config() *indicators.IndicatorConfig {
    return &indicators.IndicatorConfig{
        Name: "My Indicator",
        Parameters: []indicators.Parameter{
            {Name: "Period", DefaultValue: float64(m.period), Min: 2, Max: 100, Step: 1},
        },
        Pane: indicators.PaneSeparate,
        Outputs: []indicators.OutputConfig{
            {Name: "Value", Color: "#2196F3", Style: indicators.StyleLine, Width: 2},
        },
    }
}
```

## Key Rules

- Use ring buffers for windowed calculations (O(1) memory)
- Validate constructor inputs with `indicators.ClampMin()`
- Never allocate in `update()` or `UpdateAll()`
- Use `math.NaN()` for warmup, never `0`
- Config colors are hex strings, styles are `StyleLine`/`StyleHistogram`/`StyleDots`/`StyleArea`/`StyleDashed`
- `PaneOverlay` = drawn on price chart, `PaneSeparate` = own sub-chart
- Not thread-safe: one instance per goroutine/symbol

## Package Layout

```
types.go              — OHLCV, Indicator interface, Config types, ClampMin helpers
moving_averages/      — SMA, EMA, WMA, DEMA, TEMA, HMA, KAMA, ALMA, etc.
oscillators/          — RSI, Stochastic, CCI, Williams %R, MFI, etc.
trend/                — MACD, ADX, Parabolic SAR, Supertrend, Ichimoku, etc.
volatility/           — Bollinger, ATR, Keltner, Donchian, etc.
volume/               — OBV, VWAP, A/D Line, Chaikin MF, etc.
momentum/             — ROC, TSI, KST, Coppock, Fisher Transform, etc.
```
