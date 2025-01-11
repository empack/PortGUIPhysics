package gui

import (
	"image/color"
	"math"
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/dataDump"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/parameter"
	"physicsGUI/pkg/gui/parameter/parameter_panel"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	// App reference
	App        fyne.App
	MainWindow = dataDump.MainWindow
)

// Start GUI (function is blocking)
func Start() {
	App = app.NewWithID("GUI-Physics")
	dataDump.MainWindow = App.NewWindow("Physics GUI")

	AddMainWindow()
	App.Quit()
}

func createImportButton(window fyne.Window) *widget.Button {
	return widget.NewButton("Import Data", func() {

		// open dialog
		println("Import Data is out of function")
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

	// create dataset x^3
	dataset := make(function.Points, 21)
	for i := 0; i < len(dataset); i++ {
		dataset[i] = &function.Point{
			X:     float64(i),
			Y:     math.Pow(float64(i), 3),
			Error: 1,
		}
	}

	/* dummyGraph := graph.NewGraphCanvas(&graph.GraphConfig{
		Resolution: 100,
		Title:      "Dummy Graph to load data later",
		Function: function.NewFunction(function.Points{{
			X:     0,
			Y:     0,
			Error: 0,
		}}, function.INTERPOLATION_NONE),
	})
	GraphContainer.Add(dummyGraph) */

	profilePanel := parameter_panel.NewParameterGrid()
	for i := range data.ParameterList {
		profilePanel.Add(parameter.NewParameter(
			data.ParameterList[i].Name,
			data.ParameterList[i].DefaultVal,
			data.ParameterList[i].Val,
			data.ParameterList[i].Min,
			data.ParameterList[i].Max,
			data.ParameterList[i].Check,
		))
	}

	tickChannel := time.Tick(100 * time.Millisecond)
	go Looper(tickChannel)

	/* profilePanel.OnValueChanged = func() {
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
		sldGraph.UpdateFunction(newEdensity)
	} */

	content := container.NewBorder(
		topContainer, // top
		nil,          // bottom
		nil,          // left
		nil,          // right

		container.NewHSplit(
			container.NewVSplit(
				Graphs[0],
				profilePanel,
			),
			/* container.NewVScroll( */
			Graphs[1],
			/* ), */
		),
	)

	dataDump.MainWindow.Resize(fyne.NewSize(1000, 500))
	dataDump.MainWindow.SetContent(content)

	dataDump.MainWindow.ShowAndRun()
}

func Looper(tickChannel <-chan time.Time) {
	for {
		select {
		case <-tickChannel:
			for _, graph := range Graphs {
				graph.ReCalculate()
			}
		}
	}
}
