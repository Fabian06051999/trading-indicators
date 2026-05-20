package indicators

// ClampMin ensures a value is at least min.
func ClampMin(value, min int) int {
	if value < min {
		return min
	}
	return value
}

// ClampMinF ensures a float value is at least min.
func ClampMinF(value, min float64) float64 {
	if value < min {
		return min
	}
	return value
}

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

// PointMarker defines the shape used for StyleDots rendering in SciChart.
type PointMarker int

const (
	MarkerEllipse PointMarker = iota
	MarkerTriangle
	MarkerSquare
	MarkerCross
)

// OutputConfig describes one output series of an indicator.
type OutputConfig struct {
	Name      string
	Color     string
	UpColor   string  // conditional color when value > prev (e.g. green histogram)
	DownColor string  // conditional color when value < prev (e.g. red histogram)
	FillColor string  // fill color for Area style or band between series
	Opacity   float64 // 0.0 = transparent, 1.0 = fully opaque
	Style     LineStyle
	Width     int
	DashArray []int       // SciChart strokeDashArray, e.g. [5,5] for dashed
	Marker    PointMarker // point marker shape for StyleDots
	YRange    *YRange
	Levels    []Level
}

// IndicatorConfig holds all display metadata for an indicator.
type IndicatorConfig struct {
	Name       string
	Parameters []Parameter
	Pane       PaneType
	Outputs    []OutputConfig
}

// Indicator is the unified interface every indicator must implement.
// All indicators return [][]float64 (one slice per output series).
// Single-output indicators return one row, multi-output return multiple.
type Indicator interface {
	// CalculateAll computes all values for the given candle history (batch mode).
	CalculateAll(candles []OHLCV) [][]float64

	// UpdateAll computes the next value(s) for a single new candle (incremental mode).
	UpdateAll(candle OHLCV) []float64

	// Reset clears internal state for reuse with new data.
	Reset()

	// Config returns the display configuration with defaults.
	Config() *IndicatorConfig
}
