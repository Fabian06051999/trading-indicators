package momentum

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("ROC", "Momentum", map[string]float64{"Period": 12}, func(p map[string]float64) indicators.Indicator {
		return NewROC(int(p["Period"]))
	})
	indicators.Register("Momentum", "Momentum", map[string]float64{"Period": 10}, func(p map[string]float64) indicators.Indicator {
		return NewMomentum(int(p["Period"]))
	})
	indicators.Register("TSI", "Momentum", map[string]float64{"Long": 25, "Short": 13, "Signal": 7}, func(p map[string]float64) indicators.Indicator {
		return NewTSI(int(p["Long"]), int(p["Short"]), int(p["Signal"]))
	})
	indicators.Register("KST", "Momentum", map[string]float64{"ROC1": 10, "ROC2": 15, "ROC3": 20, "ROC4": 30, "SMA1": 10, "SMA2": 10, "SMA3": 10, "SMA4": 15, "Signal": 9}, func(p map[string]float64) indicators.Indicator {
		return NewKST(int(p["ROC1"]), int(p["ROC2"]), int(p["ROC3"]), int(p["ROC4"]), int(p["SMA1"]), int(p["SMA2"]), int(p["SMA3"]), int(p["SMA4"]), int(p["Signal"]))
	})
	indicators.Register("CoppockCurve", "Momentum", map[string]float64{"WMA": 10, "ROC1": 14, "ROC2": 11}, func(p map[string]float64) indicators.Indicator {
		return NewCoppockCurve(int(p["WMA"]), int(p["ROC1"]), int(p["ROC2"]))
	})
	indicators.Register("ElderRay", "Momentum", map[string]float64{"Period": 13}, func(p map[string]float64) indicators.Indicator {
		return NewElderRay(int(p["Period"]))
	})
	indicators.Register("WilliamsAD", "Momentum", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewWilliamsAD()
	})
	indicators.Register("ChandeForecast", "Momentum", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewChandeForecast(int(p["Period"]))
	})
	indicators.Register("PercentB", "Momentum", map[string]float64{"Period": 20, "StdDev": 2}, func(p map[string]float64) indicators.Indicator {
		return NewPercentB(int(p["Period"]), p["StdDev"])
	})
	indicators.Register("PivotPoints", "Momentum", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewPivotPoints()
	})
	indicators.Register("RVI", "Momentum", map[string]float64{"Period": 10}, func(p map[string]float64) indicators.Indicator {
		return NewRVI(int(p["Period"]))
	})
	indicators.Register("ConnorsRSI", "Momentum", map[string]float64{"RSI": 3, "Streak": 2, "Rank": 100}, func(p map[string]float64) indicators.Indicator {
		return NewConnorsRSI(int(p["RSI"]), int(p["Streak"]), int(p["Rank"]))
	})
	indicators.Register("SchaffTrendCycle", "Momentum", map[string]float64{"Fast": 23, "Slow": 50, "Cycle": 10}, func(p map[string]float64) indicators.Indicator {
		return NewSchaffTrendCycle(int(p["Fast"]), int(p["Slow"]), int(p["Cycle"]))
	})
	indicators.Register("SqueezeMomentum", "Momentum", map[string]float64{"BBPeriod": 20, "BBStdDev": 2, "KCPeriod": 20, "KCMulti": 1.5}, func(p map[string]float64) indicators.Indicator {
		return NewSqueezeMomentum(int(p["BBPeriod"]), p["BBStdDev"], int(p["KCPeriod"]), p["KCMulti"])
	})
	indicators.Register("FisherTransform", "Momentum", map[string]float64{"Period": 9}, func(p map[string]float64) indicators.Indicator {
		return NewFisherTransform(int(p["Period"]))
	})
	indicators.Register("ChoppinessIndex", "Momentum", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewChoppinessIndex(int(p["Period"]))
	})
	indicators.Register("WaveTrend", "Momentum", map[string]float64{"ChannelPeriod": 10, "AvgPeriod": 21}, func(p map[string]float64) indicators.Indicator {
		return NewWaveTrend(int(p["ChannelPeriod"]), int(p["AvgPeriod"]))
	})
}
