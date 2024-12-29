package function

import (
	"cmp"
	"math"
	"slices"
)

// represents a point with x and y value and and error value
type Point struct {
	X     float64
	Y     float64
	Error float64
}

type Points []*Point

// sorts all points by X value
func (p Points) Sort() {
	slices.SortFunc(p, func(a, b *Point) int {
		return cmp.Compare(a.X, b.X)
	})
}

// returns min and max X value of all points
func (p Points) MinMaxX() (float64, float64) {
	p.Sort()
	if len(p) == 0 {
		return 0, 0
	}
	return p[0].X, p[len(p)-1].X
}

// returns min and max Y value of all points
func (p Points) MinMaxY() (float64, float64) {
	if len(p) == 0 {
		return 0, 0
	}

	min := p[0].Y
	max := p[0].Y

	for _, point := range p {
		if max < point.Y {
			max = point.Y
		}
		if min > point.Y {
			min = point.Y
		}
	}

	return min, max
}

// ? may be useless, but may simplify some stuff so we will keep it for now
func (p Points) Log() Points {
	np := make(Points, len(p))
	for i, point := range p {
		np[i] = &Point{
			X:     point.X,
			Y:     math.Log10(point.Y),
			Error: point.Error,
		}
	}

	return p
}

// represents a coordinate with x and y value
type Coordinate struct {
	X float64
	Y float64
}
