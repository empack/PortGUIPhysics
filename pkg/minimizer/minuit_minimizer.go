package minimizer

import (
	"errors"
	"fmt"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/param"
)

type PentaltyFunction func(dtrack *MinuitFunction, parameter []float64) float64

type MinuitFunction struct {
	// function for calculating the penalty
	PenaltyFunction PentaltyFunction

	// datatracks of the experimantel data
	ExperimentalData []function.Functions

	// parameter values
	Parameters param.Parameters[float64]
}

func NewMinuitFcn(dtrack []function.Functions, pen PentaltyFunction, params param.Parameters[float64]) *MinuitFunction {
	return &MinuitFunction{
		PenaltyFunction:  pen,
		ExperimentalData: dtrack,
		Parameters:       params,
	}
}

func (d *MinuitFunction) ValueOf(par []float64) float64 {
	return d.PenaltyFunction(d, par)
}

// updates the parameters with the current values
func (d *MinuitFunction) UpdateParameters(current []float64) error {
	if len(current) != len(d.Parameters) {
		return errors.New("current values and parameters have different length")
	}

	for i, p := range d.Parameters {
		if err := p.Set(current[i]); err != nil {
			return fmt.Errorf("could not update parameter %d: %s", i, err)
		}
	}

	return nil
}
