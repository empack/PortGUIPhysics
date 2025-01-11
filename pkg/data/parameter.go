package data

import (
	"fyne.io/fyne/v2/data/binding"
	"math"
)

type Parameter struct {
	check      binding.Bool
	name       binding.String
	defaultVal binding.Float
	val        binding.Float
	min        binding.Float
	max        binding.Float
}

func NewParameter(check bool, name string, defaultVal, val, min, max float64) Parameter {
	p := Parameter{
		check:      binding.NewBool(),
		name:       binding.NewString(),
		defaultVal: binding.NewFloat(),
		val:        binding.NewFloat(),
		min:        binding.NewFloat(),
		max:        binding.NewFloat(),
	}
	_ = p.check.Set(check)
	_ = p.name.Set(name)
	_ = p.val.Set(val)
	_ = p.min.Set(min)
	_ = p.max.Set(max)
	_ = p.defaultVal.Set(defaultVal)
	return p
}

var ParameterList = [...]Parameter{
	NewParameter(false, "Eden a", 0.0, 0.0, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Eden 1", 0.346197, 0.346197, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Eden 2", 0.458849, 0.458849, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Eden b", 0.334000, 0.334000, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Roughness a/1", 3.39544, 3.39544, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Roughness 1/2", 2.15980, 2.15980, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Roughness 2/b", 3.90204, 3.90204, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Thickness 1", 14.2657, 14.2657, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "Thickness 2", 10.6906, 10.6906, -math.MaxFloat64, math.MaxFloat64),

	NewParameter(false, "deltaq", 0.0, 0.0, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "background", 0.0, 10e-9, -math.MaxFloat64, math.MaxFloat64),
	NewParameter(false, "scaling", 1.0, 1.0, -math.MaxFloat64, math.MaxFloat64),
}
