package gui

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"path/filepath"
	"physicsGUI/pkg/data"
	dataio "physicsGUI/pkg/data/io"
	"physicsGUI/pkg/data/transformation"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/graph"
	"physicsGUI/pkg/gui/parameter"
	"physicsGUI/pkg/gui/parameter/parameter_panel"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	// App reference
	App            fyne.App
	MainWindow     fyne.Window
	GraphContainer *fyne.Container
)

// Start GUI (function is blocking)
func Start() {
	App = app.NewWithID("GUI-Physics")
	MainWindow = App.NewWindow("Physics GUI")
	GraphContainer = container.NewVBox()

	AddMainWindow()
}

func createImportButton(window fyne.Window) *widget.Button {
	return widget.NewButton("Import Data", func() {

		// open dialog
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
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
				dialog.ShowError(err, window)
				return
			}

			// get filename
			filename := filepath.Base(reader.URI().Path())

			// handle import
			measurements, err := dataio.Parse(bytes)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if len(measurements) == 0 {
				dialog.ShowError(errors.New("no data"), window)
				return
			}

			points := make(function.Points, len(measurements))
			for i, m := range measurements {
				points[i] = m.ToPoint()
			}

			// convert to Point format
			/* points := make([][]function.Point, measurements[0].Count)
			for j, m := range measurements {
				for i := 0; i < measurements[j].Count; i++ {
					if j == 0 {
						points[i] = make([]function.Point, len(measurements))
					}
					points[i][j] = function.Point{
						X:     m.Time,
						Y:     m.Data[i],
						Error: m.Error,
					}
				}
			} */

			GraphContainer.RemoveAll()
			//minP, _ := plotFunc.Scope()
			plot := graph.NewGraphCanvas(&graph.GraphConfig{
				Title: fmt.Sprintf("Data track %d", 1),
				IsLog: false,
				//MinValue:   minP.X,
				Resolution: 200,
				Function:   function.NewFunction(points, function.INTERPOLATION_NONE),
			})

			GraphContainer.Add(plot)
			// Clear old plots and add new
			/* GraphContainer.RemoveAll()
			for i := 0; i < len(points); i++ {
				plotFunc := function.NewDataFunction(points, function.INTERPOLATION_NONE)
				minP, _ := plotFunc.Scope()
				plot := NewGraphCanvas(&GraphConfig{
					Title:      fmt.Sprintf("Data track %d", i+1),
					IsLog:      false,
					MinValue:   minP.X,
					Resolution: 200,
					Function:   plotFunc,
				})

				GraphContainer.Add(plot)
			}
			GraphContainer.Refresh() */

			// show success message
			dialog.ShowInformation("Import successful",
				fmt.Sprintf("File '%s' imported", filename),
				window)
		}, window)
	})
}

// separator
func createSeparator() *canvas.Line {
	line := canvas.NewLine(color.Gray{Y: 100})
	line.StrokeWidth = 1
	return line
}

// AddMainWindow builds and renders the main GUI content, it will show and run the main window,
// which is a blocking command [fyne.Window.ShowAndRun]
func AddMainWindow() {
	importButton := createImportButton(MainWindow)

	toolbar := container.NewHBox(
		importButton,
	)

	separator := createSeparator()

	topContainer := container.NewVBox(
		toolbar,
		separator,
	)

	// create dataset 2^x
	dataset := make(function.Points, 21)
	for i := 0; i < len(dataset); i++ {
		dataset[i] = &function.Point{
			X:     float64(i),
			Y:     math.Pow(float64(i), 3),
			Error: 1,
		}
	}

	g1 := graph.NewGraphCanvas(&graph.GraphConfig{
		Title:    "Non Logarithmic x³",
		IsLog:    false,
		Function: function.NewFunction(dataset, function.INTERPOLATION_NONE),
	})

	g2 := graph.NewGraphCanvas(&graph.GraphConfig{
		Title:    "Logarithmic x³",
		IsLog:    true,
		Function: function.NewFunction(dataset, function.INTERPOLATION_NONE),
	})

	GraphContainer.Add(g1)

	sldGraph := graph.NewGraphCanvas(&graph.GraphConfig{
		Resolution: 5,
		Title:      "Electron Density",
		Function:   function.NewFunction(function.Points{}, function.INTERPOLATION_NONE),
	})

	// Add all parameters to Parameter grid
	profilePanel := parameter_panel.NewParameterGrid()
	for _, param := range data.StartUpParameters {
		profilePanel.Add(parameter.NewWrapper(param))
	}

	//// HOW TO USE EXAMPLE
	// Define Data Manipulation Pipeline see data/transformation/pipeline.go
	transformationPipeline := transformation.NewBasicAsyncPipeline[[]*data.Parameter, function.Function]()
	// Set SldAdapter to handle transformation from []*data.Parameter to function.Function see function/pipeline_adapter.go
	transformationPipeline.AddStage(transformation.NewStage[[]*data.Parameter, function.Function](function.NewSldAdapter()))
	// Set SldAdapter to handle transformation from function.Function to function.Function (Graphs don't change the functions they just render it) see graph/pipeline_adapter.go
	transformationPipeline.AddStage(transformation.NewStage[function.Function, function.Function](sldGraph))

	// Setup Updater
	updater := NewScreenUpdater(transformationPipeline)
	// Start updater with 10 updates slots per seconds (updates will only be performed, when SetData was called before the update slot)
	updater.Loop(100 * time.Millisecond)
	// Update the data (The parameter will be the input for the first stage in the transformationPipeline)
	updater.SetData(data.StartUpParameters) //TODO set data again, when changed WIP-fix set updater dirty auf true use listener on parameter fields

	content := container.NewBorder(
		topContainer, // top
		nil,          // bottom
		nil,          // left
		nil,          // right

		container.NewHSplit(
			container.NewVSplit(
				sldGraph,
				profilePanel,
			),
			/* container.NewVScroll( */
			container.NewVSplit(
				g1, g2,
			),
			/* ), */
		),
	)

	MainWindow.Resize(fyne.NewSize(1000, 500))
	MainWindow.SetContent(content)

	MainWindow.ShowAndRun()
}
