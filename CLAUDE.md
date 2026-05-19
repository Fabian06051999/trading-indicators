# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

Go library providing the top 100 technical trading indicators (no drawing tools). This is a **pure calculation library** — no frontend, no API server, no UI. A colleague imports this package into his existing Go backend which feeds a SciChart frontend.

## Build & Test Commands

```bash
go build ./...
go test ./...
go test -bench=. ./...        # performance benchmarks (target: 1M+ candles)
go test -run TestSMA ./...    # run a single test
```

## Architecture

- **Delivery:** Standalone Go module (`go get`-able)
- **Input:** OHLCV candle data
- **Computation modes:** Batch (full history) + Incremental (live updates via WebSocket)
- **Performance target:** 1M+ candles — use streaming algorithms, ring buffers, pre-allocated slices, zero-copy where possible

### Indicator Interface

Every indicator implements both batch and incremental calculation, and exposes display metadata for the SciChart frontend:

```go
type Indicator interface {
    Calculate(candles []OHLCV) []float64
    Update(candle OHLCV) float64
    Reset()
    Config() *IndicatorConfig
}
```

### Display Config (for SciChart Frontend)

Each indicator carries a config struct with sensible defaults that can be modified before use:

```go
type IndicatorConfig struct {
    Name        string
    Parameters  []Parameter       // e.g. Period=14, adjustable
    Pane        PaneType          // Overlay (on price chart) or Separate (own pane)
    Outputs     []OutputConfig    // one per line/series the indicator draws
}

type OutputConfig struct {
    Name       string            // e.g. "MACD Line", "Signal"
    Color      string            // hex color default, e.g. "#2196F3"
    Style      LineStyle         // Line, Histogram, Dots, Area
    Width      int               // line thickness in px
    YRange     *YRange           // optional fixed range (e.g. 0-100 for RSI)
    Levels     []Level           // horizontal levels (e.g. 70/30 for RSI)
}

type PaneType int
const (
    PaneOverlay  PaneType = iota  // drawn on the price chart
    PaneSeparate                   // gets its own sub-chart
)
```

Design choice: Config is a separate struct returned by `Config()`. The user can modify colors, styles, parameters etc. before passing it to the frontend. Calculation logic and display config are decoupled — changing a color doesn't require recalculating.

### Package Layout

Indicators are grouped by category:
- `moving_averages/` — SMA, EMA, WMA, DEMA, TEMA, KAMA, HMA, VWMA, etc.
- `oscillators/` — RSI, Stochastic, CCI, Williams %R, MFI, CMO, etc.
- `trend/` — MACD, ADX, Parabolic SAR, Ichimoku, Supertrend, etc.
- `volatility/` — Bollinger Bands, ATR, Keltner, Donchian, Std Dev, etc.
- `volume/` — OBV, VWAP, A/D Line, Chaikin MF, etc.
- `momentum/` — ROC, Momentum, TSI, Ultimate Oscillator, etc.

Core types and interfaces live in the root package (`types.go`).

## Key Constraints

- Every indicator must work both in batch and incremental mode.
- Every indicator must expose an `IndicatorConfig` with display defaults (colors, style, pane type, levels, Y-range).
- Pane logic is simple: `PaneOverlay` (on price chart) or `PaneSeparate` (own sub-chart). SciChart handles the rest.
- Config and calculation are decoupled — modifying display properties never triggers recalculation.
- Unit tests use known reference values (e.g. from TradingView/Excel) to verify correctness.
- This repo will be handed off — all design decisions must be documented (README, doc comments on exported types).
- No frontend or API code belongs here. Only pure indicator calculation logic + display metadata.
