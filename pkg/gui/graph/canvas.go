package graph

import (
	"image/color"
	"physicsGUI/pkg/data/transformation"
	"physicsGUI/pkg/function"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// GraphCanvas represents the graphical representation of a graph.
type GraphCanvas struct {
	widget.BaseWidget
	transformation.BaseSegment[function.Function, function.Function]
	config     *GraphConfig
	background *canvas.Rectangle

	function *function.Function
}

// NewGraphCanvas creates a new canvas instance with a provided config.
// Specfically, it sets up the underlying structure of a canvas including lines, axes, labels and background.
// The method also calls 'ExtendBaseWidget' to cross-reference the canvas instance with the underlying fyne.BaseWidget struct.
func NewGraphCanvas(config *GraphConfig) *GraphCanvas {
	g := &GraphCanvas{
		config:     config,
		background: canvas.NewRectangle(color.Black),

		function: config.Function,
	}

	// needs to be to cross reference with the underlying struct
	g.ExtendBaseWidget(g)

	return g
}

// CreateRenderer returns a [GraphRenderer] from a [GraphCanvas]
func (g *GraphCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &GraphRenderer{
		graph:   g,
		objects: make([]fyne.CanvasObject, 0),
		size:    &fyne.Size{},
		margin:  float32(50),
	}
}

// UpdateFunction updates the function and refreshes the [GraphCanvas]
func (g *GraphCanvas) UpdateFunction(newFunction *function.Function) {
	g.function = newFunction
	g.Refresh()
}
