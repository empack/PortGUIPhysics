package data

import (
	"math/rand/v2"
	"slices"
	"testing"
)

const testCount = 5

func TestEdensity_len3_err1(t *testing.T) {
	result, reserr := getEdensities([]float64{0.0, 1.0, 2.0}, []float64{0.0, 1.0, 2.0}, []float64{0.0, 1.0, 2.0}, 20)

	if result != nil || reserr == nil {
		t.Fail()
	}
}
func TestEdensity_len3_1(t *testing.T) {
	for i := 0; i < testCount; i++ {
		edensitys := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		d := []float64{rand.Float64(), rand.Float64()}
		sigma := []float64{rand.Float64(), rand.Float64(), rand.Float64()}
		ref, referr := getEden(edensitys, d, sigma, 100)
		res, reserr := getEdensities(edensitys, d, sigma, 100)
		if referr != nil || reserr != nil || !slices.Equal(ref, res) {
			if referr != nil {
				t.Skipf("cant test with buged reference RefErr: %s", referr)
			} else {
				var errText = "nil"
				if reserr != nil {
					errText = reserr.Error()
				}
				t.Logf("ResError: %s,ResultEqual: %t", errText, slices.Equal(ref, res))
				t.FailNow()
			}
		}
	}
}
func TestEdensity_len7_err1(t *testing.T) {
	result, reserr := getEdensities(
		[]float64{0.0, 1.0, 2.0, 0.0, 1.0, 2.0, 42.0},
		[]float64{0.0, 1.0, 2.0, 0.0, 1.0, 2.0, 42.0},
		[]float64{0.0, 1.0, 2.0, 0.0, 1.0, 2.0, 42.0},
		100)

	if result != nil || reserr == nil {
		t.Fail()
	}
}
func TestEdensity_len7_1(t *testing.T) {
	for i := 0; i < testCount; i++ {
		edensitys := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		d := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		sigma := []float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
		ref, referr := getEden(edensitys, d, sigma, 10)
		res, reserr := getEdensities(edensitys, d, sigma, 10)
		if referr != nil || reserr != nil || !slices.Equal(ref, res) {
			if referr != nil {
				t.Skipf("cant test with buged reference RefErr: %s", referr)
			} else {
				var errText = "nil"
				if reserr != nil {
					errText = reserr.Error()
				}
				missing, additional := pointSliceDiff(ref, res)
				t.Logf("Run <%d>", i+1)
				t.Logf("ResError: %s,ResultEqual: %t", errText, slices.Equal(ref, res))
				t.Logf("Missing: %d, Additional: %d, Hitted: %d", len(missing), len(additional), len(ref)-len(missing))
				t.Fail()
			}
		}
	}
}

func pointSliceDiff(ref, res []Point) ([]Point, []Point) {
	temp := map[Point]int{}
	for _, s := range ref {
		temp[s]++
	}
	for _, s := range res {
		temp[s]--
	}

	var mis []Point
	var add []Point
	for s, v := range temp {
		if v > 0 {
			mis = append(mis, s)
		}
		if v < 0 {
			add = append(add, s)
		}
	}
	return mis, add
}
