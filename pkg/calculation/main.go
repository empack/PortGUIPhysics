package calculation

import (
	"math"
	"physicsGUI/pkg/function"
)

type Calculation struct {
	Provider  []func() function.Points
	Functions []*function.Function
}

func (this *Calculation) ReCalculate() {
	functions := make([]*function.Function, len(this.Provider))
	for i := range this.Provider {
		functions[i] = function.NewFunction(this.Provider[i](), function.INTERPOLATION_LINEAR)
	}
}

func (this *Calculation) Error(resolution int) float64 {
	diff := 0.0
	models := make([]function.Points, len(this.Provider))
	for i := range this.Functions {
		_, models[i] = this.Functions[i].Model(resolution)
	}
	for i := 0; i < resolution; i++ {
		for j := 1; j < len(models); j++ {
			diff += math.Pow(models[0][i].Y-models[j][i].Y, 2)
		}
	}
	return diff
}
