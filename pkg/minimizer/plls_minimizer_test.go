package minimizer

import (
	"math"
	"slices"
	"testing"
)

func pllsTestFunc(in []float64) float64 {
	return -1 / (in[0]*in[0] + in[1]*in[1])
}
func TestPLLS2DMinimizer(t *testing.T) {
	x0 := []float64{2, -2}
	minima := []float64{-4, -4}
	maxima := []float64{+4, +4}
	problem := NewProblem(x0, minima, maxima, pllsTestFunc, NewAsyncConfig(1e7))

	FloatMinimizerPLLS.Minimize(problem)

	res, err := problem.GetCurrentParameters()
	if err != nil {
		t.Errorf("Failed to get parameters from problem after minimizer should have finished: %s", err.Error())
	}
	if !slices.EqualFunc(res, []float64{0, 0}, func(f float64, f2 float64) bool {
		return math.Abs(f-f2) < 1e-6
	}) {
		t.Errorf("Minimizer failed to minimize value expected {0,0} but got {%f,%f}", res[0], res[1])
	}
}
