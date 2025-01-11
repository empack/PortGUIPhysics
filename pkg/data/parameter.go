package data

import (
	"fyne.io/fyne/v2/data/binding"
	"math"
)

type Parameter struct {
	Check      binding.Bool
	Name       binding.String
	DefaultVal binding.Float
	Val        binding.Float
	Min        binding.Float
	Max        binding.Float
}

func NewParameter(check bool, name string, defaultVal, val, min, max float64) Parameter {
	p := Parameter{
		Check:      binding.NewBool(),
		Name:       binding.NewString(),
		DefaultVal: binding.NewFloat(),
		Val:        binding.NewFloat(),
		Min:        binding.NewFloat(),
		Max:        binding.NewFloat(),
	}
	_ = p.Check.Set(check)
	_ = p.Name.Set(name)
	_ = p.Val.Set(val)
	_ = p.Min.Set(min)
	_ = p.Max.Set(max)
	_ = p.DefaultVal.Set(defaultVal)
	return p
}

var ParameterList = [...]Parameter{
	NewParameter(false, "Eden a", 0.0, 0.0, -math.MaxFloat64, math.MaxFloat64),                // 0
	NewParameter(false, "Eden 1", 0.346197, 0.346197, -math.MaxFloat64, math.MaxFloat64),      // 1
	NewParameter(false, "Eden 2", 0.458849, 0.458849, -math.MaxFloat64, math.MaxFloat64),      // 2
	NewParameter(false, "Eden b", 0.334000, 0.334000, -math.MaxFloat64, math.MaxFloat64),      // 3
	NewParameter(false, "Roughness a/1", 3.39544, 3.39544, -math.MaxFloat64, math.MaxFloat64), // 4
	NewParameter(false, "Roughness 1/2", 2.15980, 2.15980, -math.MaxFloat64, math.MaxFloat64), // 5
	NewParameter(false, "Roughness 2/b", 3.90204, 3.90204, -math.MaxFloat64, math.MaxFloat64), // 6
	NewParameter(false, "Thickness 1", 14.2657, 14.2657, -math.MaxFloat64, math.MaxFloat64),   // 7
	NewParameter(false, "Thickness 2", 10.6906, 10.6906, -math.MaxFloat64, math.MaxFloat64),   // 8

	NewParameter(false, "deltaq", 0.0, 0.0, -math.MaxFloat64, math.MaxFloat64),       // 9
	NewParameter(false, "background", 0.0, 10e-9, -math.MaxFloat64, math.MaxFloat64), // 10
	NewParameter(false, "scaling", 1.0, 1.0, -math.MaxFloat64, math.MaxFloat64),      // 11
}
