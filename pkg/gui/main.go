package gui

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"math/rand"
	"path/filepath"
	"physicsGUI/pkg/data"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	App        fyne.App
	MainWindow fyne.Window
)

// Start GUI (function is blocking)
func Start() {
	App = app.New()
	MainWindow = App.NewWindow("Physics GUI")

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
			defer reader.Close()

			// read file
			bytes, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			// get filename
			filename := filepath.Base(reader.URI().Path())

			// handle import
			if err := data.Import(bytes, filename); err != nil {
				dialog.ShowError(err, window)
				return
			}

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
	dataset := make([]data.Point, 10)
	for i := 0; i < 10; i++ {
		dataset[i] = data.Point{
			X:   float64(i),
			Y:   math.Pow(2, float64(i)),
			ERR: 1,
		}
	}

	graph1 := NewGraphCanvas(&GraphConfig{
		Title:    "Logarithmic",
		IsLog:    true,
		MinValue: 0.01,
		Data:     data.NewFunction(dataset, data.INTERPOLATION_NONE),
	})
	graph2 := NewGraphCanvas(&GraphConfig{
		Title: "Linear",
		Data:  data.NewFunction(dataset, data.INTERPOLATION_NONE),
	})
	graph3 := NewGraphCanvas(&GraphConfig{
		Title: "50ms Updates (bench)",
		Data:  data.NewFunction(dataset, data.INTERPOLATION_NONE),
	})
	sldGraph := NewGraphCanvas(&GraphConfig{
		Title: "SLD",
		Data:  data.NewFunction(dataset, data.INTERPOLATION_NONE),
	})

	profilePanel := NewProfilePanel(NewSldDefaultSettings("Settigns"))

	graphs := container.NewHSplit(
		graph1,
		graph2,
	)

	// "benchnmark"
	go func(graph *GraphCanvas) {
		for {
			newData := []data.Point{}
			for i := 0; i < 10; i++ {
				newData = append(newData, data.Point{
					X:   float64(i),
					Y:   float64(rand.Intn(150)),
					ERR: rand.Float64() * 20,
				})
			}

			graph.UpdateData(data.NewFunction(newData, data.INTERPOLATION_NONE))
			time.Sleep(50 * time.Millisecond)
		}
	}(graph3)

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
			container.NewVSplit(
				graphs,
				graph3,
			),
		),
	)

	MainWindow.Resize(fyne.NewSize(1000, 500))
	MainWindow.SetContent(content)

	MainWindow.ShowAndRun()
}
