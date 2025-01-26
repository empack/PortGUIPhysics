package genetic_minimizer

import (
	"slices"
	"testing"
)

var testDNAComplete = []float64{123.0, 1237895467.0, 1.9e-12, 9.4567e42}
var testDNAModifiable = []float64{1237895467.0, 1.9e-12}
var fixedDNA = []bool{true, false, false, true}

func TestDNA_GetCompletedSequence(t *testing.T) {
	dna := NewDNA(testDNAComplete, fixedDNA)
	if !slices.Equal(dna.GetCompletedSequence(), testDNAComplete) {
		t.Errorf("GetCompletedSequence returned wrong sequence")
	}
}
func TestDNA_GetModifiableSequence(t *testing.T) {
	dna := NewDNA(testDNAComplete, fixedDNA)
	if !slices.Equal(dna.GetModifiableSequence(), testDNAModifiable) {
		t.Errorf("GetModifiableSequence returned wrong sequence")
	}
}
