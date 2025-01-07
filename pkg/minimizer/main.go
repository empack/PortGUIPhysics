package minimizer

type Minimizer interface {
	Minimize(errFunc func([]float64) float64, x0 []float64, maxIterations uint, minDelta float64, acceptanceError float64) []float64
}
