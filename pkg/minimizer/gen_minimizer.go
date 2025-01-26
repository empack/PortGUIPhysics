package minimizer

import (
	"math"
	"math/rand/v2"
	"physicsGUI/pkg/minimizer/genetic_minimizer"
	"slices"
)

//// Some Default Implementations

func FitnessWeightedRandom[T Number](evalDna *[]genetic_minimizer.EvaluatedDNA[T], selectionCount int) []int {
	rate := float64(len(*evalDna))
	sectionRange := math.MaxFloat64 / float64(len(*evalDna))
	res := make([]int, selectionCount)
	index := 0
	for index < len(res) {
		val := math.Mod(rand.ExpFloat64()/rate, sectionRange)
		if !slices.Contains(res, int(val)) {
			res[index] = int(val)
			index++
		}
	}
	return res
}

func TopRated[T Number](evalDna *[]genetic_minimizer.EvaluatedDNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	for i := range res {
		res[i] = i
	}
	return res
}

func WorstRated[T Number](evalDna *[]genetic_minimizer.EvaluatedDNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	for i := range res {
		res[i] = len(*evalDna) - i
	}
	return res
}

func RandomBottomRated[T Number](evalDna *[]genetic_minimizer.EvaluatedDNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	index := 0
	for index < len(res) {
		candidate := rand.IntN(len(*evalDna)/2) + (len(*evalDna) / 2)
		if !slices.Contains(res, candidate) {
			res[index] = candidate
		}
	}
	return res
}

type GeneticMinimizer[T Number] struct {
	environment genetic_minimizer.Environment[T]
}

func (g GeneticMinimizer[T]) Minimize(problem *AsyncMinimiserProblem[T]) {

}
