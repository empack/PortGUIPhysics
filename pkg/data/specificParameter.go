package data

//var made public, beacause they needed to be used in globaly ?
var Eden_1 = 0.40
var Eden_2 = 0.30
var Eden_3 = 0.20
var Eden_b = 0.709
var D_1 = 7.0
var D_2 = 13.0
var D_3 = 5.0
var Sigma_a1 = 3.0
var Sigma_12 = 3.0
var Sigma_23 = 3.0
var Sigma_7b = 3.0
var Background = 2e-6
var Scaling = 1.0
var Eden_4 = 0.40
var D_4 = 7.0
var Sigma_34 = 3.0
var Eden_5 = 0.4
var D_5 = 3.0
var Sigma_45 = 3.0
var Eden_6 = 0.709
var D_6 = 20.0
var Sigma_56 = 3.0
var Eden_7 = 0.709
var D_7 = 20.0
var Sigma_67 = 3.0
var Deltaqz = 0.0

func TestBoundarys(parameter []float64, max []float64, min []float64) []float64 {
	if len(max) == len(min) && len(parameter) == len(max) {
		for i := 0; i < len(parameter); i++ {
			//Only positive values?
			if parameter[i] < 0 {
				parameter[i] = -parameter[i]
			}
			if parameter[i] < min[i] {
				parameter[i] = min[i]
			} else {
				if parameter[i] > max[i] {
					parameter[i] = max[i]
				}
			}
		}
	}
	return parameter
}

func TestBoundary(parameter float64, max float64, min float64) float64 {
	//Only positive values?
	if parameter < 0 {
		parameter = -parameter
	}
	if parameter < min {
		parameter = min
	} else {
		if parameter > max {
			parameter = max
		}
	}
	return parameter
}

var Parameters [26]float64

//function used to get the Array from the Sourcecode. But do we need it?
func DefineArray() {
	Parameters := [26]float64{Background, Scaling, Eden_1, Eden_2, Eden_3, Eden_b, D_1, D_2, D_3, Sigma_a1, Sigma_12, Sigma_23, Sigma_7b, Eden_4, D_4, Sigma_34, Eden_5, D_5, Sigma_45, Eden_6, D_6, Sigma_56, Eden_7, D_7, Sigma_67, Deltaqz}
	//function for checking boundarys if needed, IDL to Go-Code:
	//bounds every Parameter with specific boundarys 6 -> Parameter[6]
	var bounds = [18]int{6, 7, 8, 17, 20, 23, 9, 10, 11, 12, 15, 18, 21, 24, 2, 13, 19, 16}
	/*
	 min and max are set for stage, starting at 1 counting through the stages array.

	*/
	var stage = 0
	var max = [4]float64{30, 6, 0.34, 0.25}
	var min = [4]float64{3, 1, 0.3, 0.15}
	var stages = [4]int{6, 14, 17, 18}
	for i := 0; i < 18; i++ {
		if stages[stage] == i {
			stage++
		}
		println("Stage:", stage, " boundswert:", bounds[i])
		if Parameters[bounds[i]] < min[stage] {
			Parameters[bounds[i]] = min[stage]
		} else {
			if Parameters[bounds[i]] > max[stage] {
				Parameters[bounds[i]] = max[stage]
			}
		}
	}
	/* Old Code 1:1 from IDL:
	var bounds = [6]int64{6, 7, 8, 17, 20, 23}
	for i := 0; i < len(bounds); i++ {
		if Parameters[bounds[i]] < 3 {
			Parameters[bounds[i]] = 3
		} else {
			if Parameters[bounds[i]] > 30 {
				Parameters[bounds[i]] = 30
			}
		}
	}
	var bounds2 = [8]int{9, 10, 11, 12, 15, 18, 21, 24}
	for i := 0; i < len(bounds2); i++ {
		if Parameters[bounds2[i]] < 1 {
			Parameters[bounds2[i]] = 1
		} else {
			if Parameters[bounds2[i]] > 6 {
				Parameters[bounds2[i]] = 6
			}
		}
	}
	var bounds3 = [3]int{2, 13, 19}
	for i := 0; i < len(bounds3); i++ {
		if Parameters[bounds3[i]] < 0.3 {
			Parameters[bounds3[i]] = 0.3
		} else {
			if Parameters[bounds3[i]] > 0.34 {
				Parameters[bounds3[i]] = 0.34
			}
		}
	}
	if Parameters[16] < 0.15 {
		Parameters[16] = 0.15
	} else {
		if Parameters[16] > 0.25 {
			Parameters[16] = 0.25
		}
	}
	*/
}
