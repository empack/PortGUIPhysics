package param

import (
	"fmt"
)

type GroupElements[T any] struct {
	params []*Parameter[T]
	ref    map[string]int
}

// NewGroupElements creates a new group of elements
func NewGroupElements[T any]() *GroupElements[T] {
	return &GroupElements[T]{
		params: make([]*Parameter[T], 0),
		ref:    make(map[string]int),
	}
}

// checks if parameter is in the group
func (g *GroupElements[T]) Check(label string) bool {
	return g.ref[label] != 0
}

// adds a parameter to the group
func (g *GroupElements[T]) Add(label string, param *Parameter[T]) {
	g.params = append(g.params, param)
	g.ref[label] = len(g.params)
}

// returns the element for the specific label
func (g GroupElements[T]) GetParam(label string) *Parameter[T] {
	if element := g.params[g.ref[label]-1]; element != nil {
		return element
	}

	return nil
}

// sets the value for the specific label
func (g GroupElements[T]) Set(label string, value T) error {
	if element := g.GetParam(label); element != nil {
		return element.Set(value)
	}

	return fmt.Errorf("parameter with label %s not found", label)
}

// sets all values in the group
func (g GroupElements[T]) SetAll(values []T) error {
	if len(g.params) != len(values) {
		return fmt.Errorf("number of values does not match number of parameters in group")
	}

	// iterate over all values and set them according to the values index
	for i, value := range values {
		if err := g.params[i].Set(value); err != nil {
			return fmt.Errorf("error setting value for parameter %d: %w", i, err)
		}
	}

	return nil
}

// returns all values in the group
func (g GroupElements[T]) GetValues() ([]T, error) {
	values := make([]T, len(g.params))

	for i, param := range g.params {
		v, err := param.Get()
		if err != nil {
			return nil, fmt.Errorf("error getting value for parameter %d: %w", i, err)
		}

		values[i] = v
	}

	return values, nil
}

// returns all minima in the group
func (g GroupElements[T]) GetMinima() ([]T, error) {
	values := make([]T, len(g.params))

	for i, param := range g.params {
		if len(param.GetRelatives()) > 0 {
			v, err := param.GetRelatives()[0].Get()
			if err != nil {
				return nil, fmt.Errorf("error getting minima for parameter %d: %w", i, err)
			}

			values[i] = v
		}
	}

	return values, nil
}

// returns all maxima in the group
func (g GroupElements[T]) GetMaxima() ([]T, error) {
	values := make([]T, len(g.params))

	for i, param := range g.params {
		v, err := param.GetRelatives()[1].Get()
		if err != nil {
			return nil, fmt.Errorf("error getting value for parameter %d: %w", i, err)
		}

		values[i] = v
	}

	return values, nil
}

// returns the checkbox value of all in the group
func (g GroupElements[T]) GetChecked() []bool {
	values := make([]bool, len(g.params))

	for i, param := range g.params {
		v := param.IsChecked()

		values[i] = v
	}

	return values
}
