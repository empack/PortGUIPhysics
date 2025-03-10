package param

import (
	"errors"
	"image/color"
	"maps"
	"slices"
)

type ParameterGroup[T any] map[string]*GroupElements[T]

var (
	// sParams is a map of string parameter groups
	// each group contains a map of parameter labels and their values
	// each group can be used for iterating through parameters of the same type
	sParams = make(ParameterGroup[string])

	// fParams is a map of float parameter groups
	// each group contains a map of parameter labels and their values
	// each group can be used for iterating through parameters of the same type (edesntiy, roughness, thickness)
	fParams = make(ParameterGroup[float64])

	// iParams is a map of int parameter groups
	// each group contains a map of parameter labels and their values
	// each group can be used for iterating through parameters of the same type (limits, number of slabs)
	iParams = make(ParameterGroup[int])

	// label configs
	labelColor  = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	minMaxColor = color.NRGBA{R: 120, G: 120, B: 120, A: 255}

	// ErrParameterNotFound is returned when a parameter is not found
	ErrParameterNotFound = errors.New("parameter not found")
)

func GetStringGroup(group string) *GroupElements[string] {
	return sParams[group]
}

func GetFloatGroup(group string) *GroupElements[float64] {
	return fParams[group]
}

func GetIntGroup(group string) *GroupElements[int] {
	return iParams[group]
}
func GetStringKeys() []string {
	return slices.Collect(maps.Keys(sParams))
}
func GetFloatKeys() []string {
	return slices.Collect(maps.Keys(fParams))
}
func GetIntKeys() []string {
	return slices.Collect(maps.Keys(iParams))
}
