package io

import (
	"errors"
	"fmt"
	"math"
	"physicsGUI/pkg/function"
)

// NewOldSLDFunction returns a DataPoint function Based on the old implementation of the getEden function
//
// - eden is an array with all the eden values {eden_a,eden_1,eden_2,...,eden_n,eden_b} (edensity)
//
// - d array with the d values {d_1,d_2,...,d_n} (Thickness)
//
// - sigma array with sigma values {sigma_a1,sigma_12,sigma_23,...,sigma_(n-1)(n),sigma_nb} (Roughness)
//
// # Returns nil, when old getEden was not implemented for parameter count
//
// **NOTE** This uses old 1 to 1 translated code from IDL and has not to be defined for any parameter count
//
// Deprecated: Use new function implementation when possible
func NewOldSLDFunction(eden []float64, d []float64, sigma []float64, zNumber int) *function.Function {
	points, err := getEdensities(eden, d, sigma, zNumber)
	if err != nil {
		return nil
	}
	return function.NewFunction(points, function.INTERPOLATION_NONE)
}

var edens map[int]func(eden []float64, d []float64, sigma []float64, zNumber int) function.Points

func init() {
	edens = make(map[int]func(eden []float64, d []float64, sigma []float64, zNumber int) function.Points)
	edens[1] = eden1
	edens[2] = eden2
	edens[3] = eden3
	edens[7] = eden7
}

// created based on eden3
func eden1(edensitys []float64, d []float64, sigma []float64, znumber int) function.Points {
	zaxis := get_zaxis(d, znumber)

	eden_a := edensitys[0] //old ?
	eden_1 := edensitys[1]
	eden_b := edensitys[2]
	d_1 := d[0]
	sigma_a1 := math.Abs(sigma[0])
	sigma_1b := math.Abs(sigma[1])

	z_a1 := 0.0
	z_1b := d_1

	eden := make(function.Points, znumber)
	for i := 0; i < znumber; i++ {
		z := zaxis[i]
		step1 := (eden_1 - eden_a) * 0.5 * (1.0 + math.Erf((z-z_a1)/(math.Sqrt2*sigma_a1)))
		step2 := (eden_b - eden_1) * 0.5 * (1.0 + math.Erf((z-z_1b)/(math.Sqrt2*sigma_1b)))
		eden_i := step1 + step2
		eden[i] = &function.Point{
			X:     z,
			Y:     eden_i,
			Error: 0,
		}
	}

	return eden
}

// created based on eden3
func eden2(edensitys []float64, d []float64, sigma []float64, znumber int) function.Points {
	zaxis := get_zaxis(d, znumber)

	eden_a := edensitys[0] // old ?
	eden_1 := edensitys[1]
	eden_2 := edensitys[2]
	eden_b := edensitys[3]
	d_1 := d[0]
	d_2 := d[1]
	sigma_a1 := math.Abs(sigma[0])
	sigma_12 := math.Abs(sigma[1])
	sigma_2b := math.Abs(sigma[2])

	z_a1 := 0.0
	z_12 := d_1
	z_2b := d_1 + d_2

	eden := make(function.Points, znumber)
	for i := 0; i < znumber; i++ {
		z := zaxis[i]
		step1 := (eden_1 - eden_a) * 0.5 * (1.0 + math.Erf((z-z_a1)/(math.Sqrt2*sigma_a1)))
		step2 := (eden_2 - eden_1) * 0.5 * (1.0 + math.Erf((z-z_12)/(math.Sqrt2*sigma_12)))
		step3 := (eden_b - eden_2) * 0.5 * (1.0 + math.Erf((z-z_2b)/(math.Sqrt2*sigma_2b)))
		eden_i := step1 + step2 + step3
		eden[i] = &function.Point{
			X:     z,
			Y:     eden_i,
			Error: 0,
		}
	}

	return eden
}

// get_eden3 from fit_refl_monolayer
func eden3(edensitys []float64, d []float64, sigma []float64, znumber int) function.Points {
	zaxis := get_zaxis(d, znumber)

	eden_a := edensitys[0] //old 0.0
	eden_1 := edensitys[1]
	eden_2 := edensitys[2]
	eden_3 := edensitys[3]
	eden_b := edensitys[4]
	d_1 := d[0]
	d_2 := d[1]
	d_3 := d[2]
	sigma_a1 := math.Abs(sigma[0])
	sigma_12 := math.Abs(sigma[1])
	sigma_23 := math.Abs(sigma[2])
	sigma_3b := math.Abs(sigma[3])

	z_a1 := 0.0
	z_12 := d_1
	z_23 := d_1 + d_2
	z_3b := d_1 + d_2 + d_3

	eden := make(function.Points, znumber)
	for i := 0; i < znumber; i++ {
		z := zaxis[i]
		step1 := (eden_1 - eden_a) * 0.5 * (1.0 + math.Erf((z-z_a1)/(math.Sqrt2*sigma_a1)))
		step2 := (eden_2 - eden_1) * 0.5 * (1.0 + math.Erf((z-z_12)/(math.Sqrt2*sigma_12)))
		step3 := (eden_3 - eden_2) * 0.5 * (1.0 + math.Erf((z-z_23)/(math.Sqrt2*sigma_23)))
		step4 := (eden_b - eden_3) * 0.5 * (1.0 + math.Erf((z-z_3b)/(math.Sqrt2*sigma_3b)))
		eden_i := step1 + step2 + step3 + step4
		eden[i] = &function.Point{
			X:     z,
			Y:     eden_i,
			Error: 0,
		}
	}

	return eden
}

// get_eden from fit_bilayer_rigaku
func eden7(edensitys []float64, d []float64, sigma []float64, zNumber int) function.Points {

	zaxis := get_zaxis(d, zNumber)

	eden_a := edensitys[0] //old 0.33
	eden_1 := edensitys[1]
	eden_2 := edensitys[2]
	eden_3 := edensitys[3]
	eden_4 := edensitys[4]
	eden_5 := edensitys[5]
	eden_6 := edensitys[6]
	eden_7 := edensitys[7]
	eden_b := edensitys[8]
	d_1 := d[0]
	d_2 := d[1]
	d_3 := d[2]
	d_4 := d[3]
	d_5 := d[4]
	d_6 := d[5]
	d_7 := d[6]
	sigma_a1 := math.Abs(sigma[0])
	sigma_12 := math.Abs(sigma[1])
	sigma_23 := math.Abs(sigma[2])
	sigma_34 := math.Abs(sigma[3])
	sigma_45 := math.Abs(sigma[4])
	sigma_56 := math.Abs(sigma[5])
	sigma_67 := math.Abs(sigma[6])
	sigma_7b := math.Abs(sigma[7])

	z_a1 := 0.0
	z_12 := d_1
	z_2a3 := d_1 + d_2
	z_32b := d_1 + d_2 + d_3
	z_2b4 := d_1 + 2*d_2 + d_3
	z_45 := d_1 + 2*d_2 + d_3 + d_4
	z_56 := d_1 + 2*d_2 + d_3 + d_4 + d_5
	z_67 := d_1 + 2*d_2 + d_3 + d_4 + d_5 + d_6
	z_7b := d_1 + 2*d_2 + d_3 + d_4 + d_5 + d_6 + d_7

	eden := make(function.Points, zNumber)

	for i := 0; i < zNumber; i++ {
		z := zaxis[i]
		step1 := (eden_1 - eden_a) * 0.5 * (1.0 + math.Erf((z-z_a1)/(math.Sqrt2*sigma_a1)))
		step2 := (eden_2 - eden_1) * 0.5 * (1.0 + math.Erf((z-z_12)/(math.Sqrt2*sigma_12)))
		step3a := (eden_3 - eden_2) * 0.5 * (1.0 + math.Erf((z-z_2a3)/(math.Sqrt2*sigma_23)))
		step3b := (eden_2 - eden_3) * 0.5 * (1.0 + math.Erf((z-z_32b)/(math.Sqrt2*sigma_23)))
		step4 := (eden_4 - eden_2) * 0.5 * (1.0 + math.Erf((z-z_2b4)/(math.Sqrt2*sigma_34)))
		step5 := (eden_5 - eden_4) * 0.5 * (1.0 + math.Erf((z-z_45)/(math.Sqrt2*sigma_45)))
		step6 := (eden_6 - eden_5) * 0.5 * (1.0 + math.Erf((z-z_56)/(math.Sqrt2*sigma_56)))
		step7 := (eden_7 - eden_6) * 0.5 * (1.0 + math.Erf((z-z_67)/(math.Sqrt2*sigma_67)))
		step8 := (eden_b - eden_7) * 0.5 * (1.0 + math.Erf((z-z_7b)/(math.Sqrt2*sigma_7b)))
		eden_i := eden_a + step1 + step2 + step3a + step3b + step4 + step5 + step6 + step7 + step8
		eden[i] = &function.Point{
			X:     z,
			Y:     eden_i,
			Error: 0,
		}
	}

	return eden
}

// getEden returns a DataPoints based on the old implementation of the old getEden function
//
// - eden is an array with all the eden values {eden_a,eden_1,eden_2,...,eden_n,eden_b} (edensity)
//
// - d array with the d values {d_1,d_2,...,d_n} (Thickness)
//
// - sigma array with sigma values {sigma_a1,sigma_12,sigma_23,...,sigma_(n-1)(n),sigma_nb} (Roughness)
//
// **NOTE** This uses old 1 to 1 translated code from IDL and has not to be defined for any parameter count
//
// Deprecated: Use new function implementation when possible
func getEden(eden []float64, d []float64, sigma []float64, zNumber int) (function.Points, error) { // TODO understand this mess?
	f := edens[len(d)]
	if f == nil {
		return nil, errors.New(fmt.Sprintf("No getEden schema defined for %d edens", len(d)))
	}

	return f(eden, d, sigma, zNumber), nil
}

// getEden returns a DataPoints based on the old implementation of the old getEden function
//
// - eden is an array with all the eden values {eden_a,eden_1,eden_2,...,eden_n,eden_b} (edensity)
//
// - d array with the d values {d_1,d_2,...,d_n} (Thickness)
//
// - sigma array with sigma values {sigma_a1,sigma_12,sigma_23,...,sigma_(n-1)(n),sigma_nb} (Roughness)
func getEdensities(eden []float64, d []float64, sigma []float64, zNumber int) ([]*function.Point, error) {

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

	edensities := make([]*function.Point, zNumber)
	zaxis := get_zaxis(d, zNumber)

	for i := 0; i < zNumber; i++ {
		z_i := zaxis[i]
		//calculate cumulative edensity at a specific z_i
		y := 0.0
		for step := 0; step < step_n; step++ {
			y += (eden[step+1] - eden[step]) * 0.5 * (1.0 + math.Erf((z_i-z[step])/(math.Sqrt2*math.Abs(sigma[step]))))
		}
		//create points for drawing
		edensities[i] = &function.Point{
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
