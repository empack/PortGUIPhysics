package genetic_minimizer

import (
	"math/rand"
	"slices"
)

type DNA[T Number] struct {
	modifiableSequence []T
	completedTemplate  []T
	fixedMap           []bool
}

func NewDNA[T Number](gens []T, fixed []bool) *DNA[T] {
	falseCnt := 0
	for _, b := range fixed {
		if !b {
			falseCnt++
		}
	}
	completedTemplate := slices.Clone(gens)
	fixedMap := slices.Clone(fixed)
	modifiableSequence := make([]T, falseCnt)
	modCnt := 0
	for i := range gens {
		if !fixed[i] {
			modifiableSequence[modCnt] = gens[i]
			modCnt++
		}
	}
	return &DNA[T]{
		modifiableSequence: modifiableSequence,
		completedTemplate:  completedTemplate,
		fixedMap:           fixedMap,
	}
}

func (dna *DNA[T]) Mutate(rate, amplifier float64) {
	for i := range dna.modifiableSequence {
		if rand.Float64() > rate {
			dna.modifiableSequence[i] += T(rand.Float64() * amplifier)
		}
	}
}

func (dna *DNA[T]) GetModifiableSequence() []T {
	return dna.modifiableSequence
}

func (dna *DNA[T]) GetCompletedSequence() []T {
	res := slices.Clone(dna.completedTemplate)
	modCnt := 0
	for i := range res {
		if !dna.fixedMap[i] {
			res[i] = dna.modifiableSequence[modCnt]
			modCnt++
		}
	}
	return res
}
