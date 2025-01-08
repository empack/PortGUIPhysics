package minimizer

type Minimizer interface {
	Minimize(errFunc func([]float64) float64, x0 []float64, maxIterations int, minDelta float64) []float64
}

var HillClimbingMinimizer = hillClimbingMinimizer{}
