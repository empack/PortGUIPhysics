package graph

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"path/filepath"
	"physicsGUI/pkg/dataDump"
	"physicsGUI/pkg/function"
	mod_io "physicsGUI/pkg/io"
)

// GraphCanvas represents the graphical representation of a graph.
type GraphCanvas struct {
	widget.BaseWidget
	config        *GraphConfig
	background    *canvas.Rectangle
	btnImportData *widget.Button
	functions     []*function.Function
}

// NewGraphCanvas creates a new canvas instance with a provided config.
// Specfically, it sets up the underlying structure of a canvas including lines, axes, labels and background.
// The method also calls 'ExtendBaseWidget' to cross-reference the canvas instance with the underlying fyne.BaseWidget struct.
func NewGraphCanvas(config *GraphConfig) *GraphCanvas {
	g := &GraphCanvas{
		config:        config,
		background:    canvas.NewRectangle(color.Black),
		btnImportData: nil,
		functions:     config.Functions,
	}
	g.btnImportData = widget.NewButton("📲", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, dataDump.MainWindow)
				return
			}
			if reader == nil {
				return // user canceled
			}
			defer func() {
				if err := reader.Close(); err != nil {
					log.Println("error while closing dialog:", err)
				}
			}()

			// read file
			bytes, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, dataDump.MainWindow)
				return
			}

			// get filename
			filename := filepath.Base(reader.URI().Path())

			// handle import
			measurements, err := mod_io.Parse(bytes)
			if err != nil {
				dialog.ShowError(err, dataDump.MainWindow)
				return
			}

			if len(measurements) == 0 {
				dialog.ShowError(errors.New("no data"), dataDump.MainWindow)
				return
			}

			points := make(function.Points, len(measurements))
			for i, m := range measurements {
				points[i] = m.ToPoint()
			}
			g.config.Calculation = append(g.config.Calculation, func() function.Points {
				return points
			})
			g.ReCalculate()

			// show success message
			dialog.ShowInformation("Import successful",
				fmt.Sprintf("File '%s' imported", filename),
				dataDump.MainWindow)
		}, dataDump.MainWindow)

	})

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
func (g *GraphCanvas) UpdateFunction(newFunctions []*function.Function) {
	g.functions = newFunctions
	g.Refresh()
}

func (r *GraphCanvas) ReCalculate() {
	functions := make([]*function.Function, len(r.config.Calculation))
	for i := range r.config.Calculation {
		functions[i] = function.NewFunction(r.config.Calculation[i](), function.INTERPOLATION_NONE)
	}
	r.UpdateFunction(functions)
}
