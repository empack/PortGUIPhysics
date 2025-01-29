package minimizer

import (
	"math"
	"math/rand/v2"
	"physicsGUI/pkg/minimizer/genetic_minimizer"
	"slices"
)

//// Some Default SelectionBehavior Implementations

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

//// Some Default CrossoverBehavior Implementations

func TwoParentSinglePointCrossover[T Number](dna ...*genetic_minimizer.DNA[T]) {
	s1 := dna[0].GetModifiableSequence()
	s2 := dna[1].GetModifiableSequence()
	splitIndex := rand.IntN(len(s1))
	s1a := s1[:splitIndex]
	s2a := s2[splitIndex:]
	s1b := s1[splitIndex:]
	s2b := s2[:splitIndex]
	dna[0].SetModifiableSequence(append(s1a, s2a...))
	dna[1].SetModifiableSequence(append(s1b, s2b...))
}

func TwoParentTwoPointCrossover[T Number](dna ...*genetic_minimizer.DNA[T]) {
	s1 := dna[0].GetModifiableSequence()
	s2 := dna[1].GetModifiableSequence()
	splitIndex1 := rand.IntN(len(s1))
	splitIndex2 := rand.IntN(len(s1)-splitIndex1) + splitIndex1
	s1a := s1[:splitIndex1]
	s2a := s2[:splitIndex1]
	s1b := s1[splitIndex1:splitIndex2]
	s2b := s2[splitIndex1:splitIndex2]
	s1c := s1[splitIndex2:]
	s2c := s2[splitIndex2:]
	dna[0].SetModifiableSequence(append(s1a, append(s2b, s1c...)...))
	dna[1].SetModifiableSequence(append(s2a, append(s1b, s2c...)...))
}

func ThreeParentCrossover[T Number](dna ...*genetic_minimizer.DNA[T]) {

}

type GeneticMinimizer[T Number] struct {
	environment genetic_minimizer.Environment[T]
}

func (g GeneticMinimizer[T]) Minimize(problem *AsyncMinimiserProblem[T]) {

}
