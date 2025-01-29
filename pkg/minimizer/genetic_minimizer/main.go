package genetic_minimizer

type Number interface {
	~uint8 | ~uint32 | ~uint64 |
		~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}
type EvaluatedDNA[T Number] struct {
	float64
	*DNA[T]
}
type FitnessEvaluator[T Number] func(dna *DNA[T]) float64
type CrossoverBehavior[T Number] func(dna ...*DNA[T])
type SelectionBehavior[T Number] func(evalDna *[]EvaluatedDNA[T], selectionCount int) []int

type Environment[T Number] struct {
	MutationRate       float64
	MutationAmplifier  float64
	PopulationSize     int
	WildcardCount      int
	PrecursorCount     int
	ParentCount        int
	FitnessEvaluator   FitnessEvaluator[T]
	CrossoverBehavior  CrossoverBehavior[T]
	ParentSelection    SelectionBehavior[T]
	WildcardSelection  SelectionBehavior[T]
	PrecursorSelection SelectionBehavior[T]
}
