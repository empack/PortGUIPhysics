package graph

import (
	"image/color"
	"physicsGUI/pkg/function"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// GraphCanvas represents the graphical representation of a graph.
type GraphCanvas struct {
	widget.BaseWidget
	config     *GraphConfig
	background *canvas.Rectangle

	functions function.Functions
}

// NewGraphCanvas creates a new canvas instance with a provided config.
// Specfically, it sets up the underlying structure of a canvas including lines, axes, labels and background.
// The method also calls 'ExtendBaseWidget' to cross-reference the canvas instance with the underlying fyne.BaseWidget struct.
func NewGraphCanvas(config *GraphConfig) *GraphCanvas {
	g := &GraphCanvas{
		config:     config,
		background: canvas.NewRectangle(color.Black),

		functions: config.Functions,
	}

	for _, f := range g.functions {
		if f == nil {
			panic("function cannot be nil. Make sure to provide a function (even an empty one)")
		}
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

// UpdateFunctions updates the function and refreshes the [GraphCanvas]
func (g *GraphCanvas) UpdateFunctions(newFunctions function.Functions) {
	g.functions = newFunctions
	g.Refresh()
}

func (g *GraphCanvas) AddFunction(newFunction *function.Function) {
	g.functions = append(g.functions, newFunction)
	g.Refresh()
}
