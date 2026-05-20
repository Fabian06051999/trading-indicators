package indicators

import "fmt"

// IndicatorFactory creates an indicator instance from parameters.
type IndicatorFactory func(params map[string]float64) Indicator

// Registration holds metadata for a registered indicator.
type Registration struct {
	Name     string
	Category string
	Factory  IndicatorFactory
	Defaults map[string]float64
}

var registry = map[string]Registration{}

// Register adds an indicator to the global registry.
func Register(name, category string, defaults map[string]float64, factory IndicatorFactory) {
	registry[name] = Registration{
		Name:     name,
		Category: category,
		Factory:  factory,
		Defaults: defaults,
	}
}

// Create instantiates an indicator by name with the given parameters.
// Missing parameters are filled with defaults.
func Create(name string, params map[string]float64) (Indicator, error) {
	reg, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("indicator %q not found", name)
	}

	merged := make(map[string]float64)
	for k, v := range reg.Defaults {
		merged[k] = v
	}
	for k, v := range params {
		merged[k] = v
	}

	return reg.Factory(merged), nil
}

// MustCreate is like Create but panics on error.
func MustCreate(name string, params map[string]float64) Indicator {
	ind, err := Create(name, params)
	if err != nil {
		panic(err)
	}
	return ind
}

// ListAll returns all registered indicators.
func ListAll() []Registration {
	list := make([]Registration, 0, len(registry))
	for _, r := range registry {
		list = append(list, r)
	}
	return list
}

// ListCategory returns all indicators in a category.
func ListCategory(category string) []Registration {
	var list []Registration
	for _, r := range registry {
		if r.Category == category {
			list = append(list, r)
		}
	}
	return list
}

// Names returns all registered indicator names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
