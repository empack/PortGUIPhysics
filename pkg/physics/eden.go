package physics

import (
	"math"
	"physicsGUI/pkg/function"
)

//You need some different edensity calculation?
// => copy this function until you have the number of variants you need
// => give each function a unique name (GetEdensities1/GetXYEdensities)
// => insert the kind of parameters the calculation needs to the first bracket GetEdensities(param_1 type_1, ..., param_n type_n)
// => insert your calculation
// => continue by adapting RecalculateData and penaltyFunction in PortGUIPhysics\pkg\gui\main.go

// GetEdensities returns DataPoints based on the old implementation of the old getEden function
// - eden is an array with all the eden values {eden_a,eden_1,eden_2,...,eden_n,eden_b} (edensity)
// - d array with the d values {d_1,d_2,...,d_n} (Thickness)
// - sigma array with sigma values {sigma_a1,sigma_12,sigma_23,...,sigma_(n-1)(n),sigma_nb} (Roughness)
func GetEdensities(eden, size []float64, roughness, coverage float64) (function.Points, error) {

	eden_au := eden[0]
	eden_org := eden[1]
	eden_b := eden[2]

	radius := size[0]
	d_shell := size[1]
	z_offset := size[2]

	edensities := make(function.Points, ZNUMBER)
	volfrac_au := make([]float64, ZNUMBER)
	volfrac_org := make([]float64, ZNUMBER)
	volfrac_w := make([]float64, ZNUMBER)
	area := math.Pi * math.Pow(radius+d_shell, 2)

	zaxis := GetZAxis(radius, d_shell, ZNUMBER)

	for i := 0; i < ZNUMBER; i++ {
		z_i := zaxis[i]
		step_a_w := 0.5 * (1.0 + math.Erf((z_i-z_offset)/(math.Sqrt2*roughness)))
		value_au := 0.0
		value_org := 0.0
		if math.Abs(z_i) < (radius + d_shell) {
			if math.Abs(z_i) <= radius {
				value_au = math.Pi * (math.Pow(radius, 2) - math.Pow(z_i, 2)) / area
				value_org = math.Pi * ((math.Pow(radius+d_shell, 2) - math.Pow(z_i, 2)) - (math.Pow(radius, 2) - math.Pow(z_i, 2))) / area
			} else {
				value_org = math.Pi * (math.Pow(radius+d_shell, 2) - math.Pow(z_i, 2)) / area
			}
		}
		volfrac_au[i] = coverage * value_au
		volfrac_org[i] = coverage * value_org
		volfrac_w[i] = (coverage*(1.0-(value_au+value_org)) + (1.0 - coverage)) * step_a_w
		y := eden_au*volfrac_au[i] + eden_org*volfrac_org[i] + eden_b*volfrac_w[i]

		//create points for drawing
		edensities[i] = &function.Point{
			X:     z_i,
			Y:     y,
			Error: 0,
		}
	}

	return edensities, nil
}

func GetZAxis(radius, d_shell float64, zNumber int) []float64 {
	z0 := -(radius + d_shell + 30.0)
	z1 := -z0

	zStep := (z1 - z0) / float64(zNumber)
	zAxis := make([]float64, zNumber)
	for i := 0; i < zNumber; i++ {
		zAxis[i] = z0 + float64(i)*zStep
	}

	return zAxis
}
