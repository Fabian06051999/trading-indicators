package moving_averages

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("SMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewSMA(int(p["Period"]))
	})
	indicators.Register("EMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewEMA(int(p["Period"]))
	})
	indicators.Register("WMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewWMA(int(p["Period"]))
	})
	indicators.Register("DEMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewDEMA(int(p["Period"]))
	})
	indicators.Register("TEMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewTEMA(int(p["Period"]))
	})
	indicators.Register("HMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewHMA(int(p["Period"]))
	})
	indicators.Register("VWMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewVWMA(int(p["Period"]))
	})
	indicators.Register("SMMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewSMMA(int(p["Period"]))
	})
	indicators.Register("KAMA", "Moving Averages", map[string]float64{"Period": 10, "Fast": 2, "Slow": 30}, func(p map[string]float64) indicators.Indicator {
		return NewKAMA(int(p["Period"]), int(p["Fast"]), int(p["Slow"]))
	})
	indicators.Register("ZLEMA", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewZLEMA(int(p["Period"]))
	})
	indicators.Register("ALMA", "Moving Averages", map[string]float64{"Period": 9, "Offset": 0.85, "Sigma": 6}, func(p map[string]float64) indicators.Indicator {
		return NewALMA(int(p["Period"]), p["Offset"], p["Sigma"])
	})
	indicators.Register("LSMA", "Moving Averages", map[string]float64{"Period": 25}, func(p map[string]float64) indicators.Indicator {
		return NewLSMA(int(p["Period"]))
	})
	indicators.Register("McGinley", "Moving Averages", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewMcGinleyDynamic(int(p["Period"]))
	})
	indicators.Register("T3", "Moving Averages", map[string]float64{"Period": 5, "VolumeFactor": 0.7}, func(p map[string]float64) indicators.Indicator {
		return NewT3(int(p["Period"]), p["VolumeFactor"])
	})
	indicators.Register("VIDYA", "Moving Averages", map[string]float64{"Period": 14, "CMOPeriod": 9}, func(p map[string]float64) indicators.Indicator {
		return NewVIDYA(int(p["Period"]), int(p["CMOPeriod"]))
	})
	indicators.Register("FRAMA", "Moving Averages", map[string]float64{"Period": 16}, func(p map[string]float64) indicators.Indicator {
		return NewFRAMA(int(p["Period"]))
	})
}
