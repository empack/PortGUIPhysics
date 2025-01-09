package graph

import (
	"image/color"
	"physicsGUI/pkg/function"
)

var (
	// color of the title
	titleColor = &color.White

	// color of the axes
	axesColor = &color.White

	// color of the legend at the axis
	legendColor = &color.White

	// colors for the graph
	lineColor = &color.NRGBA{R: 0, G: 0, B: 255, A: 255}

	// colors for the points
	pointColor = &color.NRGBA{R: 0, G: 255, B: 0, A: 255}

	// colors for the error
	errorColor = &color.NRGBA{R: 255, G: 0, B: 0, A: 128}

	// color for the gridlines
	gridColor = &color.NRGBA{R: 128, G: 128, B: 128, A: 64}

	// color for the minor gridlines (for log scale)
	gridMinorColor = &color.NRGBA{R: 128, G: 128, B: 128, A: 48}

	// size of the points
	pointRadius = float32(0.5)
)

// GraphConfig configures the basic struct for a graph
type GraphConfig struct {
	Title      string
	IsLog      bool
	Resolution int
	Functions  []*function.Function
}
