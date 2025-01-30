package minimizer

import (
	"cmp"
	"math"
	"math/rand/v2"
	"physicsGUI/pkg/minimizer/genetic_minimizer"
	"slices"
	"sync"
)

//// Some Default SelectionBehavior Implementations

func FitnessWeightedRandom[T Number](evalDna []*genetic_minimizer.DNA[T], selectionCount int) []int {
	rate := float64(len(evalDna))
	sectionRange := math.MaxFloat64 / float64(len(evalDna))
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

func TopRated[T Number](evalDna []*genetic_minimizer.DNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	for i := range res {
		res[i] = i
	}
	return res
}

func WorstRated[T Number](evalDna []*genetic_minimizer.DNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	for i := range res {
		res[i] = len(evalDna) - i
	}
	return res
}

func RandomBottomRated[T Number](evalDna []*genetic_minimizer.DNA[T], selectionCount int) []int {
	res := make([]int, selectionCount)
	index := 0
	for index < len(res) {
		candidate := rand.IntN(len(evalDna)/2) + (len(evalDna) / 2)
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
	// TODO
	panic("implement me")
}

// random helper // TODO
func randomTBetween[T Number](min T, max T) T {
	return T(rand.Float64()*(float64(max-min)) + float64(min))
}

// helper functions for setup
func spreadOverGlobe[T Number](fillCount int, min []T, max []T, fixed []bool) []*genetic_minimizer.DNA[T] {
	dnas := make([]*genetic_minimizer.DNA[T], fillCount)

	for i := 0; i < fillCount; i++ {
		params := make([]T, len(fixed))
		for p := 0; p < len(params); p++ {
			params[p] = randomTBetween(min[p], max[p])
		}
		dnas[i] = genetic_minimizer.NewDNA(params, fixed)
	}

	return dnas
}

type GeneticMinimizer[T Number] struct {
	environment genetic_minimizer.Environment[T]
}

func (g *GeneticMinimizer[T]) Minimize(problem *AsyncMinimiserProblem[T]) {
	// Set fitness evaluator to negative error function because low error means high fitness
	g.environment.FitnessEvaluator = func(dna *genetic_minimizer.DNA[T]) float64 {
		return -problem.errorFunction(dna.GetCompletedSequence())
	}

	// create DNA from current solution
	dna0 := genetic_minimizer.NewDNA(problem.parameter, problem.fixed)
	// create random spread over min max space
	initialSeed := spreadOverGlobe(g.environment.PopulationSize-1, problem.minima, problem.maxima, problem.fixed)
	// setup first generation
	generation := append(initialSeed, dna0)

	problem.lock.RLock()
	maxGenerationCount := problem.config.LoopCount
	problem.lock.RUnlock()

	for i := 0; i < maxGenerationCount; i++ {
		// calculate fitness of generation
		wg := &sync.WaitGroup{}
		wg.Add(len(generation))
		for gen := 0; gen < len(generation); gen++ {
			go func(i int) {
				generation[i].Fitness = g.environment.FitnessEvaluator(generation[i])
				wg.Done()
			}(gen)
		}

		// sort based on fitness
		slices.SortFunc(generation, func(a, b *genetic_minimizer.DNA[T]) int {
			return cmp.Compare(a.Fitness, b.Fitness)
		})

		// fill new generation with selections
		newGeneration := make([]*genetic_minimizer.DNA[T], len(generation))

		wildcardChoice := g.environment.WildcardSelection(generation, g.environment.WildcardCount)
		precourserChoice := g.environment.PrecursorSelection(generation, g.environment.PrecursorCount)
		// TODO paralleisieren in go routinen
		for c := 0; c < len(wildcardChoice); c++ {
			newGeneration[c] = generation[wildcardChoice[c]].Clone()
		}
		for c := 0; c < len(precourserChoice); c++ {
			newGeneration[len(wildcardChoice)+c] = generation[precourserChoice[c]].Clone()
		}
		for c := g.environment.WildcardCount + g.environment.PrecursorCount; c < g.environment.PopulationSize; c += g.environment.ParentCount {
			parents := g.environment.ParentSelection(generation, g.environment.ParentCount)
			parentAdr := make([]*genetic_minimizer.DNA[T], len(parents))
			for p := 0; p < len(parents); p++ {
				newGeneration[c+p] = generation[parents[p]].Clone()
				parentAdr[p] = newGeneration[c+p]
			}
			g.environment.CrossoverBehavior(parentAdr...)
		}

		// apply mutations
		for c := 0; c < len(newGeneration); c++ {
			newGeneration[c].Mutate(g.environment.MutationRate, g.environment.MutationAmplifier)
		}

		// update problem
		problem.lock.Lock()
		problem.config.LoopCount -= 1
		problem.parameter = generation[0].GetCompletedSequence()
		problem.lock.Unlock()

		// set new generation
		generation = newGeneration
	}

	//TODO
	panic("implement me")
}
