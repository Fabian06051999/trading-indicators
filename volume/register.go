package volume

import "github.com/Fabian06051999/trading-indicators"

func init() {
	indicators.Register("OBV", "Volume", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewOBV()
	})
	indicators.Register("VWAP", "Volume", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewVWAP()
	})
	indicators.Register("ADLine", "Volume", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewADLine()
	})
	indicators.Register("ChaikinMF", "Volume", map[string]float64{"Period": 20}, func(p map[string]float64) indicators.Indicator {
		return NewChaikinMF(int(p["Period"]))
	})
	indicators.Register("ForceIndex", "Volume", map[string]float64{"Period": 13}, func(p map[string]float64) indicators.Indicator {
		return NewForceIndex(int(p["Period"]))
	})
	indicators.Register("VolumeOscillator", "Volume", map[string]float64{"Fast": 5, "Slow": 20}, func(p map[string]float64) indicators.Indicator {
		return NewVolumeOscillator(int(p["Fast"]), int(p["Slow"]))
	})
	indicators.Register("EaseOfMovement", "Volume", map[string]float64{"Period": 14}, func(p map[string]float64) indicators.Indicator {
		return NewEaseOfMovement(int(p["Period"]))
	})
	indicators.Register("NVI", "Volume", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewNVI()
	})
	indicators.Register("PVI", "Volume", map[string]float64{}, func(p map[string]float64) indicators.Indicator {
		return NewPVI()
	})
	indicators.Register("VolumeProfile", "Volume", map[string]float64{"Period": 50, "Bins": 24}, func(p map[string]float64) indicators.Indicator {
		return NewVolumeProfile(int(p["Period"]), int(p["Bins"]))
	})
}
