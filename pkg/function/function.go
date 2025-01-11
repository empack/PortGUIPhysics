package function

import (
	"errors"
	"fmt"
	"math"
)

// function with interpolation capabilities
type Function struct {
	data  Points
	Scope *Scope
	eval  InterpolationFunction
}

// scope of a function
type Scope struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

// returns a new function with the given data and interpolation mode
func NewFunction(data Points, interpolationMode InterpolationMode) *Function {
	// get min, max values of function
	minX, maxX := data.MinMaxX()
	minY, maxY := data.MinMaxY()

	// create function
	f := &Function{
		data: data,
		Scope: &Scope{
			minX,
			maxX,
			minY,
			maxY,
		},
	}

	// sanitize data
	f.Sanitize()

	// set interpolation function
	f.SetInterpolation(interpolationMode)

	return f
}

// TODO: add full explanation
func (f *Function) Model(resolution int) (Points, Points) {
	if f.eval == nil {
		return f.data, f.data
	}

	interpolated := make(Points, resolution)

	deltaX := (f.Scope.MaxX - f.Scope.MinX) / float64(resolution)
	for i := range interpolated {
		x := f.Scope.MinX + float64(i)*deltaX

		// TODO: add error handling
		y, _ := f.eval(f.data, x)

		interpolated[i] = &Point{
			X:     x,
			Y:     y,
			Error: 0,
		}
	}

	return f.data, interpolated
}

// evaluates function value at x
func (f *Function) Eval(x float64) (float64, error) {
	return f.eval(f.data, x)
}

// sanitizes the function data and removes all 0 values for potential log scale issues
// TODO: add point y handling back, but for now we only need x value handling
func (f *Function) Sanitize() {
	for i, point := range f.data {
		if point.X == 0 /* || point.Y == 0  */ {
			f.data = append(f.data[:i], f.data[i+1:]...)
		}
	}
}

// sets the interpolation function
func (f *Function) SetInterpolation(interpolationMode InterpolationMode) {
	switch interpolationMode {
	case INTERPOLATION_NONE:
		break
	case INTERPOLATION_LINEAR:
		f.eval = linearInterpolation
		break
	case INTERPOLATION_SPLINE:
	case INTERPOLATION_PCHIP:
		panic("SetInterpolation error: interpolation not implemented yet")
	default:
		panic("SetInterpolation error: Unknown interpolationMode. Please use only values provided by related const's in data package")
	}
}

// TODO: add full explanation
type FunctionSegment struct {
	start float64
	end   float64
	minY  float64
	maxY  float64
	f     *func(x float64) float64
}

func NewFunctionSegment(start float64, end float64, f *func(x float64) float64) FunctionSegment {
	//TODO Use finder for extrema?
	// find zero of difference?
	strct := FunctionSegment{
		start: start,
		end:   end,
		minY:  0,
		maxY:  0,
		f:     f,
	}

	// BEGIN TMP use end and start, assuming it monotone grow
	y1 := (*f)(start)
	y2 := (*f)(end)
	if y1 < y2 {
		strct.minY = y1
		strct.maxY = y2
	} else {
		strct.minY = y2
		strct.maxY = y1
	}
	// END TMP

	return strct
}

type SegmentedFunction struct {
	segments []FunctionSegment
	minY     float64
	maxY     float64
}

func NewSegmentedFunction(segments []FunctionSegment) *SegmentedFunction {
	var minY float64 = 0
	var maxY float64 = 0

	if len(segments) != 0 {
		var segBorder = segments[0].start
		for i, segment := range segments {
			if segBorder != segment.start {
				println(fmt.Sprintf("FunctionSegment_Warning: Gap between befinition between segments %d and %d function isn't defined contunius!", i-1, i))
			}
			if segment.minY < minY {
				minY = segment.minY
			}
			if segment.maxY > maxY {
				maxY = segment.maxY
			}
			segBorder = segment.end
		}
	}

	return &SegmentedFunction{
		segments: segments,
		minY:     minY,
		maxY:     maxY,
	}
}

func (s *SegmentedFunction) Scope() (*Coordinate, *Coordinate) {
	if s.segments == nil || len(s.segments) == 0 {
		return &Coordinate{}, &Coordinate{}
	}

	return &Coordinate{s.segments[0].start, s.minY}, &Coordinate{s.segments[len(s.segments)-1].end, s.maxY}
}

func (this *SegmentedFunction) Model(resolution int) ([]Point, []Point) {
	if this.segments == nil || len(this.segments) == 0 {
		return nil, nil
	}
	dx := (this.segments[len(this.segments)-1].end - this.segments[0].start) / float64(resolution)
	res := make([]Point, resolution)
	var nx = this.segments[0].start
	for i := 0; i < resolution; i++ {
		val, err := this.Eval(nx)
		if err != nil {
			println(err.Error())
			res[i] = Point{
				X:     0,
				Y:     0,
				Error: 0,
			}
		}
		res[i] = Point{
			X:     nx,
			Y:     val,
			Error: 0,
		}
		nx += dx
	}
	return nil, res
}

func (this *SegmentedFunction) Eval(x float64) (float64, error) {
	if this.segments == nil || len(this.segments) == 0 {
		return 0, errors.New("Evaluation_Error: Unable to evaluate when no segments are defined")
	}
	if x < this.segments[0].start || x > this.segments[len(this.segments)-1].end {
		return 0, errors.New(fmt.Sprintf("Evaluation_Error: Out of bounds %f is not in the scoupe between %f and %f", x, this.segments[0].start, this.segments[len(this.segments)-1].end))
	}
	for i := 0; i < len(this.segments); i++ {
		if this.segments[i].end >= x {
			return (*this.segments[i].f)(x), nil
		}
	}
	panic("Evaluation_Error: This should not be able to happen")
}

type MathFunctionProvider interface {
	GetF() func(x float64) float64
}
type LogisticFunction struct {
	leftBorder  float64
	growRate    float64
	rightBorder float64
	offsetX     float64
}

func NewLogisticFunction(offsetX float64, leftBorder float64, growRate float64, rightBorder float64) LogisticFunction {
	return LogisticFunction{
		leftBorder:  leftBorder,
		growRate:    growRate,
		rightBorder: rightBorder,
		offsetX:     offsetX,
	}
}
func (this LogisticFunction) GetF() func(x float64) float64 {
	return func(x float64) float64 {
		return (this.rightBorder-this.leftBorder)/(1+math.Pow(math.E, -this.growRate*(x-this.offsetX))) + this.leftBorder
	}
}
