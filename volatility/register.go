package volatility

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("BollingerBands", "Volatility", map[string]float64{"Period": 20, "StdDev": 2}, func(p map[string]float64) indicators.Indicator {
		return NewBollingerBands(int(p["Period"]), p["StdDev"])
	})
	indicators.Register("ATR", "Volatility", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewATR(int(p["Period"]))
	})
	indicators.Register("ATRPercent", "Volatility", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewATRPercent(int(p["Period"]))
	})
	indicators.Register("KeltnerChannel", "Volatility", map[string]float64{"EMAPeriod": 20, "ATRPeriod": 10, "Multiplier": 1.5}, func(p map[string]float64) indicators.Indicator {
		return NewKeltnerChannel(int(p["EMAPeriod"]), int(p["ATRPeriod"]), p["Multiplier"])
	})
	indicators.Register("DonchianChannel", "Volatility", map[string]float64{"Period": 20}, func(p map[string]float64) indicators.Indicator {
		return NewDonchianChannel(int(p["Period"]))
	})
	indicators.Register("StdDev", "Volatility", map[string]float64{"Period": 20}, func(p map[string]float64) indicators.Indicator {
		return NewStdDev(int(p["Period"]))
	})
	indicators.Register("HistoricalVolatility", "Volatility", map[string]float64{"Period": 20, "Annualization": 252}, func(p map[string]float64) indicators.Indicator {
		return NewHistoricalVolatility(int(p["Period"]), p["Annualization"])
	})
	indicators.Register("ChaikinVolatility", "Volatility", map[string]float64{"EMAPeriod": 10, "ROCPeriod": 10}, func(p map[string]float64) indicators.Indicator {
		return NewChaikinVolatility(int(p["EMAPeriod"]), int(p["ROCPeriod"]))
	})
	indicators.Register("UlcerIndex", "Volatility", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewUlcerIndex(int(p["Period"]))
	})
}
