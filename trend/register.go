package trend

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("MACD", "Trend", map[string]float64{"Fast": 12, "Slow": 26, "Signal": 9}, func(p map[string]float64) indicators.Indicator {
		return NewMACD(int(p["Fast"]), int(p["Slow"]), int(p["Signal"]))
	})
	indicators.Register("ADX", "Trend", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewADX(int(p["Period"]))
	})
	indicators.Register("ParabolicSAR", "Trend", map[string]float64{"Step": 0.02, "Max": 0.2}, func(p map[string]float64) indicators.Indicator {
		return NewParabolicSAR(p["Step"], p["Max"])
	})
	indicators.Register("Supertrend", "Trend", map[string]float64{"Period": 10, "Multiplier": 3}, func(p map[string]float64) indicators.Indicator {
		return NewSupertrend(int(p["Period"]), p["Multiplier"])
	})
	indicators.Register("Aroon", "Trend", map[string]float64{"Period": 25}, func(p map[string]float64) indicators.Indicator {
		return NewAroon(int(p["Period"]))
	})
	indicators.Register("Vortex", "Trend", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewVortex(int(p["Period"]))
	})
	indicators.Register("Ichimoku", "Trend", map[string]float64{"Tenkan": 9, "Kijun": 26, "SenkouB": 52, "Displacement": 26}, func(p map[string]float64) indicators.Indicator {
		return NewIchimoku(int(p["Tenkan"]), int(p["Kijun"]), int(p["SenkouB"]), int(p["Displacement"]))
	})
	indicators.Register("TRIX", "Trend", map[string]float64{"Period": 15}, func(p map[string]float64) indicators.Indicator {
		return NewTRIX(int(p["Period"]))
	})
	indicators.Register("DPO", "Trend", map[string]float64{"Period": 20}, func(p map[string]float64) indicators.Indicator {
		return NewDPO(int(p["Period"]))
	})
	indicators.Register("MassIndex", "Trend", map[string]float64{"EMAPeriod": 9, "SumPeriod": 25}, func(p map[string]float64) indicators.Indicator {
		return NewMassIndex(int(p["EMAPeriod"]), int(p["SumPeriod"]))
	})
	indicators.Register("EMAEnvelope", "Trend", map[string]float64{"Period": 20, "Percentage": 2.5}, func(p map[string]float64) indicators.Indicator {
		return NewEMAEnvelope(int(p["Period"]), p["Percentage"])
	})
}
