package indicators

// OHLCV represents a single candlestick with timestamp.
type OHLCV struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// PaneType defines where the indicator is rendered in SciChart.
type PaneType int

const (
	PaneOverlay  PaneType = iota // drawn on the price chart
	PaneSeparate                 // gets its own sub-chart
)

// LineStyle defines how an output series is rendered.
type LineStyle int

const (
	StyleLine      LineStyle = iota // solid line
	StyleHistogram                  // vertical bars
	StyleDots                       // scatter points
	StyleArea                       // filled area
	StyleDashed                     // dashed line
	StyleStepLine                   // step/staircase line
)

// YRange defines a fixed Y-axis range for an indicator output.
type YRange struct {
	Min float64
	Max float64
}

// Level defines a horizontal reference line (e.g. overbought/oversold).
type Level struct {
	Value float64
	Label string
	Color string
}

// Parameter defines a configurable input parameter for an indicator.
type Parameter struct {
	Name         string
	DefaultValue float64
	Min          float64
	Max          float64
	Step         float64
}

// OutputConfig describes one output series of an indicator.
type OutputConfig struct {
	Name   string
	Color  string
	Style  LineStyle
	Width  int
	YRange *YRange
	Levels []Level
}

// IndicatorConfig holds all display metadata for an indicator.
type IndicatorConfig struct {
	Name       string
	Parameters []Parameter
	Pane       PaneType
	Outputs    []OutputConfig
}

// Indicator is the core interface every indicator must implement.
type Indicator interface {
	// Calculate computes all values for the given candle history (batch mode).
	Calculate(candles []OHLCV) []float64

	// Update computes the next value for a single new candle (incremental mode).
	Update(candle OHLCV) float64

	// Reset clears internal state for reuse with new data.
	Reset()

	// Config returns the display configuration with defaults.
	Config() *IndicatorConfig
}

// MultiOutputIndicator is for indicators that produce multiple series (e.g. MACD, Bollinger).
type MultiOutputIndicator interface {
	CalculateAll(candles []OHLCV) [][]float64
	UpdateAll(candle OHLCV) []float64
	Reset()
	Config() *IndicatorConfig
}
