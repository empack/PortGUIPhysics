package function

import (
	"errors"
	"fmt"
	"math"
	"physicsGUI/pkg/data"
)

type SldFunction struct {
	Function
	eden    []*data.Parameter
	d       []*data.Parameter
	sigma   []*data.Parameter
	zNumber int
}

func (f *SldFunction) DataChanged() {
	f.updateFunction()
}

func NewSLDFunction(eden []*data.Parameter, d []*data.Parameter, sigma []*data.Parameter, zNumber int) *SldFunction {
	f := &SldFunction{
		eden:    eden,
		d:       d,
		sigma:   sigma,
		zNumber: zNumber,
	}
	for _, parameter := range eden {
		parameter.DataChannel.AddListener(f)
	}
	for _, parameter := range d {
		parameter.DataChannel.AddListener(f)
	}
	for _, parameter := range sigma {
		parameter.DataChannel.AddListener(f)
	}
	f.updateFunction()
	return f
}

func (f *SldFunction) updateFunction() {
	edens := make([]float64, len(f.eden))
	for i := range f.eden {
		if val, err := f.eden[i].GetValue().Get(); err == nil {
			edens[i] = val
		} else {
			edens[i] = 0
		}
	}
	d := make([]float64, len(f.d))
	for i := range f.d {
		if val, err := f.d[i].GetValue().Get(); err == nil {
			d[i] = val
		} else {
			d[i] = 0
		}
	}
	sigma := make([]float64, len(f.sigma))
	for i := range f.sigma {
		if val, err := f.sigma[i].GetValue().Get(); err == nil {
			sigma[i] = val
		} else {
			sigma[i] = 0
		}
	}
	if edensities, err := getEdensities(edens, d, sigma, f.zNumber); err == nil {
		newF := NewFunction(edensities, INTERPOLATION_LINEAR)
		f.Function = *newF
	}
}

// getEden returns a DataPoints based on the old implementation of the old getEden function
//
// - eden is an array with all the eden values {eden_a,eden_1,eden_2,...,eden_n,eden_b} (edensity)
//
// - d array with the d values {d_1,d_2,...,d_n} (Thickness)
//
// - sigma array with sigma values {sigma_a1,sigma_12,sigma_23,...,sigma_(n-1)(n),sigma_nb} (Roughness)
func getEdensities(eden []float64, d []float64, sigma []float64, zNumber int) ([]*Point, error) {

	step_n := len(d) + 1
	//throw error if the param number does not match the scheme
	if len(eden) != step_n+1 {
		return nil, errors.New(fmt.Sprintf("Missmatch in parameter dimensionality edensities %d/thickness %d", len(eden), len(d)))
	}
	if len(sigma) != step_n {
		return nil, errors.New(fmt.Sprintf("Missmatch in parameter dimensionality roughness %d/thickness %d", len(sigma), len(d)))
	}

	//calculate distances
	var z = make([]float64, step_n)
	z[0] = 0.0
	for i := 1; i < step_n; i++ {
		z[i] = z[i-1] + d[i-1]
	}

	edensities := make([]*Point, zNumber)
	zaxis := get_zaxis(d, zNumber)

	for i := 0; i < zNumber; i++ {
		z_i := zaxis[i]
		//calculate cumulative edensity at a specific z_i
		y := 0.0
		for step := 0; step < step_n; step++ {
			y += (eden[step+1] - eden[step]) * 0.5 * (1.0 + math.Erf((z_i-z[step])/(math.Sqrt2*math.Abs(sigma[step]))))
		}
		//create points for drawing
		edensities[i] = &Point{
			X:     z_i,
			Y:     y,
			Error: 0,
		}
	}

	return edensities, nil
}

func get_zaxis(d []float64, zNumber int) []float64 {
	z0 := -20.0
	var z1 = 0.0
	if len(d) > 3 {
		z1 += d[1]
	}
	for _, f := range d {
		z1 += f
	}
	z1 -= z0
	zStep := (z1 - z0) / float64(zNumber)
	zAxis := make([]float64, zNumber)
	for i := 0; i < zNumber; i++ {
		zAxis[i] = z0 + float64(i)*zStep
	}
	return zAxis
}
