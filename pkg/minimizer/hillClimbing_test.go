package minimizer

import (
	"math"
	"slices"
	"testing"
)

func testFunc(in []float64) float64 {
	return -1 / (in[0]*in[0] + in[1]*in[1])
}
func TestHillClimb2DMinimizer(t *testing.T) {
	minStepSize := 1e-7
	x0 := []float64{4, 4}
	minimizer := HillClimbingMinimizer
	res := minimizer.Minimize(testFunc, x0, 1e7, minStepSize)
	if !slices.EqualFunc(res, []float64{0, 0}, func(f float64, f2 float64) bool {
		return math.Abs(f-f2) < minStepSize
	}) {
		t.Errorf("Minimizer failed to minimize value expected {0,0} but got {%f,%f}", res[0], res[1])
	}
}
