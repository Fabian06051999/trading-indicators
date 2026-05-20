package oscillators

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("RSI", "Oscillators", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewRSI(int(p["Period"]))
	})
	indicators.Register("Stochastic", "Oscillators", map[string]float64{"K": 14, "D": 3, "Slowing": 3}, func(p map[string]float64) indicators.Indicator {
		return NewStochastic(int(p["K"]), int(p["D"]), int(p["Slowing"]))
	})
	indicators.Register("StochRSI", "Oscillators", map[string]float64{"RSIPeriod": 14, "K": 14, "D": 3}, func(p map[string]float64) indicators.Indicator {
		return NewStochRSI(int(p["RSIPeriod"]), int(p["K"]), int(p["D"]))
	})
	indicators.Register("CCI", "Oscillators", map[string]float64{"Period": 20}, func(p map[string]float64) indicators.Indicator {
		return NewCCI(int(p["Period"]))
	})
	indicators.Register("WilliamsR", "Oscillators", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewWilliamsR(int(p["Period"]))
	})
	indicators.Register("MFI", "Oscillators", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewMFI(int(p["Period"]))
	})
	indicators.Register("CMO", "Oscillators", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewCMO(int(p["Period"]))
	})
	indicators.Register("DeMarker", "Oscillators", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewDeMarker(int(p["Period"]))
	})
	indicators.Register("AwesomeOscillator", "Oscillators", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewAwesomeOscillator()
	})
	indicators.Register("AcceleratorOscillator", "Oscillators", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewAcceleratorOscillator()
	})
	indicators.Register("UltimateOscillator", "Oscillators", map[string]float64{"Period1": 7, "Period2": 14, "Period3": 28}, func(p map[string]float64) indicators.Indicator {
		return NewUltimateOscillator(int(p["Period1"]), int(p["Period2"]), int(p["Period3"]))
	})
	indicators.Register("PPO", "Oscillators", map[string]float64{"Fast": 12, "Slow": 26, "Signal": 9}, func(p map[string]float64) indicators.Indicator {
		return NewPPO(int(p["Fast"]), int(p["Slow"]), int(p["Signal"]))
	})
}
