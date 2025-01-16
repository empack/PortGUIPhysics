package gui

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/graph"
	"physicsGUI/pkg/gui/helper"
	"physicsGUI/pkg/gui/param"
	"physicsGUI/pkg/minimizer"
	"physicsGUI/pkg/physics"
	"physicsGUI/pkg/trigger"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	// App reference
	App        fyne.App
	MainWindow fyne.Window

	functionMap = make(map[string]*function.Function)
	graphMap    = make(map[string]*graph.GraphCanvas)
)

// Start GUI (function is blocking)
func Start() {
	App = app.NewWithID("GUI-Physics")
	App.Settings().SetTheme(theme.DarkTheme()) //TODO WIP to fix invisable while parameter lables
	MainWindow = App.NewWindow("Physics GUI")

	mainWindow()
}

// parses a given file into a dataset
func addDataset(reader io.ReadCloser, uri fyne.URI, err error) function.Points {
	if err != nil {
		dialog.ShowError(err, MainWindow)
		return nil
	}
	if reader == nil {
		return nil // user canceled
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Println("error while closing dialog:", err)
		}
	}()

	// read file
	bytes, err := io.ReadAll(reader)
	if err != nil {
		dialog.ShowError(err, MainWindow)
		return nil
	}

	// get filename
	filename := filepath.Base(uri.Name())

	// handle import
	points, err := data.Parse(bytes)
	if err != nil {
		dialog.ShowError(err, MainWindow)
		return nil
	}

	if len(points) == 0 {
		dialog.ShowError(errors.New("no data"), MainWindow)
		return nil
	}

	// show success message
	dialog.ShowInformation("Import successful",
		fmt.Sprintf("File '%s' imported", filename),
		MainWindow)

	return points
}

func createMinimizeButton() *widget.Button {
	return widget.NewButton("Minimize", func() {
		if problem := createMinimizerProblem(); problem != nil {
			closeChan := make(chan struct{}, 1)
			clock := time.Tick(1 * time.Second)

			go minimizeRefreshWorker(problem, closeChan, clock)

			go func() {
				minimizer.FloatMinimizerHC.Minimize(problem)
				closeChan <- struct{}{}
				dialog.ShowInformation("Minimizer", "Minimisation completed",
					MainWindow)
			}()

			log.Println("Minimization started")
			return
		}

		dialog.ShowInformation("Minimizer", "Minimisation Failed: Failed to create minimization Problem", MainWindow)
	})
}

func minimizeRefreshWorker(problem *minimizer.AsyncMinimiserProblem[float64], close <-chan struct{}, clock <-chan time.Time) {
	for {
		select {
		case <-close:
			return
		case <-clock:
			parameters, err := problem.GetCurrentParameters()
			if err != nil {
				continue
			}

			param.SetFloats("eden", []float64{parameters[0], parameters[1], parameters[2], parameters[3]})
			param.SetFloats("thick", []float64{parameters[4], parameters[5]})
			param.SetFloats("rough", []float64{parameters[6], parameters[7], parameters[8]})

			param.SetFloat("general", "deltaq", parameters[9])
			param.SetFloat("general", "background", parameters[10])
			param.SetFloat("general", "scaling", parameters[11])
		}
	}
}

func createMinimizerProblem() *minimizer.AsyncMinimiserProblem[float64] {
	eden, err := param.GetFloats("eden")
	if err != nil {
		return nil
	}
	d, err := param.GetFloats("thick")
	if err != nil {
		return nil
	}
	sigma, err := param.GetFloats("rough")
	if err != nil {
		return nil
	}

	// get general parameters
	delta, err := param.GetFloat("general", "deltaq")
	if err != nil {
		return nil
	}
	background, err := param.GetFloat("general", "background")
	if err != nil {
		return nil
	}
	scaling, err := param.GetFloat("general", "scaling")
	if err != nil {
		return nil
	}

	parameters := slices.Concat(eden, d, sigma)
	parameters = append(parameters, delta, background, scaling)

	dataTracks := graphMap["sld"].GetDataTracks()
	if len(dataTracks) == 0 {
		return nil
	}
	dataTrack := dataTracks[0]

	// Define error function
	errorFunction := func(params []float64) float64 {
		edenErr := params[0:4]
		dErr := params[4:6]
		sigmaErr := params[6:9]
		deltaErr := params[9]
		backgroundErr := params[10]
		scalingErr := params[11]

		edenPoints, err := physics.GetEdensities(edenErr, dErr, sigmaErr)
		if err != nil {
			fmt.Println("Error while calculating edensities:", err)
			return math.MaxFloat64
		}

		intensityPoints := physics.CalculateIntensityPoints(edenPoints, deltaErr, &physics.IntensityOptions{
			Background: backgroundErr,
			Scaling:    scalingErr,
		})

		intensityFunction := function.NewFunction(intensityPoints, function.INTERPOLATION_LINEAR)

		dataModel, _ := dataTrack.Model(0, false)

		diff := 0.0
		for i := range dataModel {
			iy, err := intensityFunction.Eval(dataModel[i].X)
			if err != nil {
				fmt.Println("Error while calculating intensity:", err)
			}
			diff += math.Abs(dataModel[i].Y - iy)
		}
		return diff
	}

	minima := make([]float64, len(parameters))
	maxima := make([]float64, len(parameters))

	for i := range minima {
		minima[i] = -math.MaxFloat64
		maxima[i] = math.MaxFloat64
	}

	return minimizer.NewProblem(parameters, minima, maxima, errorFunction, &minimizer.MinimiserConfig{
		LoopCount:     1e6,
		ParallelReads: true,
	})
}

// register functions which can be used for graph plotting
func registerFunctions() {
	functionMap["sld"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
	functionMap["eden"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
	functionMap["test"] = function.NewFunction(function.Points{}, function.INTERPOLATION_NONE)
}

// creates the graph containers for the different graphs
func registerGraphs() *fyne.Container {
	graphMap["sld"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Intensity Graph",
		IsLog:     true,
		Functions: function.Functions{functionMap["sld"]},
		DisplayRange: &graph.GraphRange{
			Min: 0.01,
			Max: math.MaxFloat64,
		},
	})

	graphMap["eden"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Edensity Graph",
		IsLog:     false,
		Functions: function.Functions{functionMap["eden"]},
	})

	graphMap["test"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Test Graph",
		IsLog:     true,
		Functions: function.Functions{functionMap["test"]},
	})

	return container.NewGridWithColumns(2, graphMap["eden"], graphMap["sld"], graphMap["test"])
}

// creates and registers the parameter and adds them to the parameter repository
func registerParams() *fyne.Container {
	edenA, _ := param.FloatMinMax("eden", "Eden a", 0.0)
	eden1, _ := param.FloatMinMax("eden", "Eden 1", 0.346197)
	eden2, _ := param.FloatMinMax("eden", "Eden 2", 0.458849)
	edenB, _ := param.FloatMinMax("eden", "Eden b", 0.334000)

	roughnessA1, _ := param.FloatMinMax("rough", "Roughness a/1", 3.39544)
	roughness12, _ := param.FloatMinMax("rough", "Roughness 1/2", 2.15980)
	roughness2B, _ := param.FloatMinMax("rough", "Roughness 2/b", 3.90204)

	thickness1, _ := param.FloatMinMax("thick", "Thickness 1", 14.2657)
	thickness2, _ := param.FloatMinMax("thick", "Thickness 2", 10.6906)

	deltaQ, _ := param.Float("general", "deltaq", 0.0)
	background, _ := param.Float("general", "background", 10e-9)
	scaling, _ := param.Float("general", "scaling", 1.0)

	return container.NewVBox(
		container.NewGridWithColumns(4, edenA, eden1, eden2, edenB),
		container.NewGridWithColumns(4, roughnessA1, roughness12, roughness2B),
		container.NewGridWithColumns(4, thickness1, thickness2),
		container.NewGridWithColumns(4, deltaQ, background, scaling),
	)
}

// onDrop is called when a file is dropped into the window and triggers
// an import of the data if the file is dropped on a graph canvas
func onDrop(position fyne.Position, uri []fyne.URI) {
	for mapIdentifier, u := range graphMap {
		if u.MouseInCanvas(position) {
			for _, v := range uri {
				rc, err := os.OpenFile(v.Path(), os.O_RDONLY, 0666)
				if err != nil {
					dialog.ShowError(err, MainWindow)
					return
				}

				if points := addDataset(rc, v, nil); points != nil {
					points.Magie()
					newFunction := function.NewFunction(points, function.INTERPOLATION_NONE)
					graphMap[mapIdentifier].AddDataTrack(newFunction)
				}
			}
			return
		}
	}
}

// mainWindow builds and renders the main GUI content, it will show and run the main window,
// which is a blocking command [fyne.Window.ShowAndRun]
func mainWindow() {
	registerFunctions()

	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				createMinimizeButton(),
			),
			helper.CreateSeparator(),
		), // top
		nil, // bottom
		nil, // left
		nil, // right

		container.NewVSplit(
			registerGraphs(),
			registerParams(),
		),
	)

	// set onchange function for recalculating data
	trigger.SetOnChange(RecalculateData)

	MainWindow.Resize(fyne.NewSize(1000, 500))
	MainWindow.SetContent(content)
	MainWindow.SetOnDropped(onDrop)

	MainWindow.ShowAndRun()
}

// this test func will create a basic x^2 dataset for testing
// and set it to the sld and eden graphs
func testFunc() {
	counter := 11

	d := make(function.Points, counter)

	for i := 0; i < counter; i++ {
		d[i] = &function.Point{
			X:     float64(i),
			Y:     math.Pow(float64(i), 2),
			Error: 1,
		}
	}

	functionMap["sld"].SetData(d)
	functionMap["eden"].SetData(d)
}

// RecalculateData recalculates the data for the sld and eden graphs
func RecalculateData() {
	// Get current parameter groups
	eden, err := param.GetFloats("eden")
	if err != nil {
		log.Println("Error while getting eden parameters:", err)
		return
	}
	d, err := param.GetFloats("thick")
	if err != nil {
		log.Println("Error while getting thickness parameters:", err)
		return
	}
	sigma, err := param.GetFloats("rough")
	if err != nil {
		log.Println("Error while getting roughness parameters:", err)
		return
	}

	// get general parameters
	delta, err := param.GetFloat("general", "deltaq")
	if err != nil {
		log.Println("Error while getting deltaq parameter:", err)
		return
	}
	background, err := param.GetFloat("general", "background")
	if err != nil {
		log.Println("Error while getting background parameter:", err)
		return
	}
	scaling, err := param.GetFloat("general", "scaling")
	if err != nil {
		log.Println("Error while getting scaling parameter:", err)
		return
	}

	// calculate edensity
	edenPoints, err := physics.GetEdensities(eden, d, sigma)
	if err != nil {
		log.Println("Error while calculating edensities:", err)
		return
	} else {
		functionMap["eden"].SetData(edenPoints)
	}

	intensityPoints := physics.CalculateIntensityPoints(edenPoints, delta, &physics.IntensityOptions{
		Background: background,
		Scaling:    scaling,
	})

	functionMap["sld"].SetData(intensityPoints)
}
