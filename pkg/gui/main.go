package gui

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"math"
	"path/filepath"
	"physicsGUI/pkg/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
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
			measurements, importErr := data.Import(bytes, filename)
			if importErr != nil {
				dialog.ShowError(err, window)
				return
			}
			if len(measurements) == 0 {
				dialog.ShowError(fmt.Errorf("no data"), window)
				return
			}

			// convert to Point format
			points := make([][]data.Point, measurements[0].Count)
			for j, m := range measurements {
				for i := 0; i < measurements[j].Count; i++ {
					if j == 0 {
						points[i] = make([]data.Point, len(measurements))
					}
					points[i][j] = data.Point{
						X:   m.Time,
						Y:   m.Data[i],
						ERR: m.Error,
					}
				}
			}

			// Clear old plots and add new
			GraphContainer.RemoveAll()
			for i := 0; i < len(points); i++ {
				plotFunc := data.NewDataFunction(points[i], data.INTERPOLATION_NONE)
				minP, _ := plotFunc.Scope()
				plot := NewGraphCanvas(&GraphConfig{
					Title:      fmt.Sprintf("Data track %d", i+1),
					IsLog:      false,
					MinValue:   minP.X,
					Resolution: 200,
					Data:       plotFunc,
				})

				GraphContainer.Add(plot)
			}
			GraphContainer.Refresh()

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

	/*
		graph1 := NewGraphCanvas(&GraphConfig{
			Title:    "Logarithmic",
			IsLog:    true,
			MinValue: 0.01,
			Data:     data.NewDataFunction(dataset, data.INTERPOLATION_NONE),
		})
	*/
	dummyFunction := data.NewOldSLDFunction(
		[]float64{0.0, 0.346197, 0.458849, 0.334000},
		[]float64{14.2657, 10.6906},
		[]float64{3.39544, 2.15980, 3.90204},
		150) // from refl_monolayer.pro:780
	if dummyFunction == nil {
		dummyFunction = data.NewDataFunction([]data.Point{{
			X:   0,
			Y:   0,
			ERR: 0,
		}}, data.INTERPOLATION_NONE)
	}
	sldGraph := NewGraphCanvas(&GraphConfig{
		Resolution: 100,
		Title:      "Electron Density ",
		Data:       dummyFunction,
	})

	dummyGraph := NewGraphCanvas(&GraphConfig{
		Resolution: 100,
		Title:      "Dummy Graph to load data later",
		Data: data.NewDataFunction([]data.Point{{
			X:   0,
			Y:   0,
			ERR: 0,
		}}, data.INTERPOLATION_NONE),
	})
	GraphContainer.Add(dummyGraph)

	profilePanel := NewProfilePanel(NewSldDefaultSettings("Settings"))
	profilePanel.OnValueChanged = func() {
		edensity := make([]float64, len(profilePanel.Profiles)+2)
		sigma := make([]float64, len(profilePanel.Profiles)+1)
		d := make([]float64, len(profilePanel.Profiles))

		var err error = nil
		edensity[0], err = profilePanel.base.Parameter[ProfileDefaultEdensityID].GetValue()
		sigma[0], err = profilePanel.base.Parameter[ProfileDefaultRoughnessID].GetValue()
		edensity[len(profilePanel.Profiles)+1], err = profilePanel.bulk.Parameter[ProfileDefaultEdensityID].GetValue()
		for i, profile := range profilePanel.Profiles {
			edensity[i+1], err = profile.Parameter[ProfileDefaultEdensityID].GetValue()
			sigma[i+1], err = profile.Parameter[ProfileDefaultRoughnessID].GetValue()
			d[i], err = profile.Parameter[ProfileDefaultThicknessID].GetValue()
		}
		var zNumberF float64 = 100.0
		zNumberF, err = profilePanel.sldSettings.Parameter[SldDefaultZNumberID].GetValue()
		zNumber := int(zNumberF)

		if err != nil {
			println(errors.Join(errors.New("error while reading default parameters"), err).Error())
		}

		newEdensity := data.NewOldSLDFunction(edensity, d, sigma, zNumber)
		if newEdensity == nil {
			println(errors.New("no old getEden function implemented for this parameter count").Error())
			return
		}
		sldGraph.UpdateData(newEdensity)
	}

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
			container.NewVScroll(GraphContainer),
		),
	)

	MainWindow.Resize(fyne.NewSize(1000, 500))
	MainWindow.SetContent(content)

	MainWindow.ShowAndRun()
}
