package gui

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"math/rand"
	"path/filepath"
	"physicsGUI/pkg/data"
	"strconv"
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
		Data:     data.NewDataFunction(dataset, data.INTERPOLATION_NONE),
	})
	graph2 := NewGraphCanvas(&GraphConfig{
		Title: "Linear",
		Data:  data.NewDataFunction(dataset, data.INTERPOLATION_NONE),
	})
	graph3 := NewGraphCanvas(&GraphConfig{
		Title: "50ms Updates (bench)",
		Data:  data.NewDataFunction(dataset, data.INTERPOLATION_NONE),
	})
	sldGraph := NewGraphCanvas(&GraphConfig{
		Title: "SLD",
		Data:  data.NewDataFunction(dataset, data.INTERPOLATION_NONE),
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

			graph.UpdateData(data.NewDataFunction(newData, data.INTERPOLATION_NONE))
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
func DefineButtons() [9]*Profile {

	p0 := NewDefaultProfile("Specific Parameters", data.ParametersName[0], data.Parameters[0], data.ParametersName[4], data.Parameters[4], data.ParametersName[7], data.Parameters[7])
	p1 := NewDefaultProfile("Specific Parameters", data.ParametersName[1], data.Parameters[1], data.ParametersName[5], data.Parameters[5], data.ParametersName[8], data.Parameters[8])
	p2 := NewDefaultProfile("Specific Parameters", data.ParametersName[2], data.Parameters[2], data.ParametersName[6], data.Parameters[6], data.ParametersName[9], data.Parameters[9])
	p3 := NewDefaultProfile("Specific Parameters", data.ParametersName[13], data.Parameters[13], data.ParametersName[14], data.Parameters[14], data.ParametersName[15], data.Parameters[15])
	p4 := NewDefaultProfile("Specific Parameters", data.ParametersName[16], data.Parameters[16], data.ParametersName[17], data.Parameters[17], data.ParametersName[18], data.Parameters[18])
	p5 := NewDefaultProfile("Specific Parameters", data.ParametersName[19], data.Parameters[19], data.ParametersName[20], data.Parameters[20], data.ParametersName[21], data.Parameters[21])
	p6 := NewDefaultProfile("Specific Parameters", data.ParametersName[22], data.Parameters[22], data.ParametersName[23], data.Parameters[23], data.ParametersName[24], data.Parameters[24])
	p7 := NewDefaultProfile("Specific Parameters", data.ParametersName[3], data.Parameters[3], data.ParametersName[25], data.Parameters[25], data.ParametersName[10], data.Parameters[10])
	p8 := NewDefaultProfile("Specific Parameters", data.ParametersName[11], data.Parameters[11], data.ParametersName[12], data.Parameters[12], "null", 0)
	proprofil := [9]*Profile{p0, p1, p2, p3, p4, p5, p6, p7, p8}
	for i := 0; i < 8; i++ {
		proprofil[i] = NewDefaultProfile(("Specific Parameters" + strconv.Itoa(i)), data.ParametersName[3*i], data.Parameters[3*i], data.ParametersName[3*i+1], data.Parameters[3*i+1], data.ParametersName[3*i+2], data.Parameters[3*i+2])
	}
	proprofil[8] = NewDefaultProfile("Specific Parameters"+"8", data.ParametersName[24], data.Parameters[24], data.ParametersName[25], data.Parameters[25], "Nameless", 1)
	//proprofil := NewDefaultProfile("Specific Parameters", data.ParametersName[0], data.Parameters[0], data.ParametersName[1], data.Parameters[1], data.ParametersName[2], data.Parameters[2])
	/*for i := 3; i < len(data.Parameters); i++ {
		proprofil.AddParameter(NewParameter(data.ParametersName[i], data.Parameters[i]))
		println(data.ParametersName[i], data.Parameters[i])
	}
	print(len(proprofil.parameter))
	*/
	return proprofil
}
