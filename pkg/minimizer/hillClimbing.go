package minimizer

type HillClimbingMinimizer struct {
}

func (h HillClimbingMinimizer) Minimize(errFunc func([]float64) float64, x0 []float64, maxIterations int, minDelta float64, acceptanceError float64) []float64 {
	guess := x0
	for i := 0; i < maxIterations; i++ {
		err := errFunc(guess)
		if err <= acceptanceError {
			return guess
		}

	}
	return guess
}
