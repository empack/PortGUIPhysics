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
		errorStr, problem := createMinimizerProblem()

		if problem != nil {
			closeChan := make(chan struct{}, 1)
			clock := time.Tick(1 * time.Second)

			go minimizeRefreshWorker(problem, closeChan, clock)

			go func() {
				minimizer.FloatMinimizerGen.Minimize(problem)
				closeChan <- struct{}{}
				dialog.ShowInformation("Minimizer", "Minimisation completed",
					MainWindow)
			}()

			log.Println("Minimization started")
			return
		}

		dialogFailedText := fmt.Sprintf("Failed to create minimization problem with error: %s", errorStr)

		dialog.ShowInformation("Minimizer", dialogFailedText, MainWindow)
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

func createMinimizerProblem() (string, *minimizer.AsyncMinimiserProblem[float64]) {
	eden, err := param.GetFloats("eden")
	if err != nil {
		return "EDEN FLOATS MISS", nil
	}
	edenMinima, err := param.GetFloatMinima("eden")
	if err != nil {
		return "EDEN MINIMA MISS", nil
	}
	edenMaxima, err := param.GetFloatMaximas("eden")
	if err != nil {
		return "EDEN MAXIMA MISS", nil
	}

	d, err := param.GetFloats("thick")
	if err != nil {
		return "THICK FLOATS MISS", nil
	}
	dMinima, err := param.GetFloatMinima("thick")
	if err != nil {
		return "THICK MINIMA MISS", nil
	}
	dMaxima, err := param.GetFloatMaximas("thick")
	if err != nil {
		return "THICK MAXIMA MISS", nil
	}

	sigma, err := param.GetFloats("rough")
	if err != nil {
		return "ROUGH FLOATS MISS", nil
	}
	sigmaMinima, err := param.GetFloatMinima("rough")
	if err != nil {
		return "ROUGH MINIMA MISS", nil
	}
	sigmaMaxima, err := param.GetFloatMaximas("rough")
	if err != nil {
		return "ROUGH MAXIMA MISS", nil
	}

	// get general parameters
	delta, err := param.GetFloat("general", "deltaq")
	if err != nil {
		return "(GENERAL) DELTAQ MISS", nil
	}
	deltaMinima := -1.0
	deltaMaxima := 1.0

	background, err := param.GetFloat("general", "background")
	if err != nil {
		return "(GENERAL) BACKGROUND MISS", nil
	}
	backgroundMinima := 0.0
	backgroundMaxima := 10.0

	scaling, err := param.GetFloat("general", "scaling")
	if err != nil {
		return "(GENERAL) SCALING MISS", nil
	}
	scalingMinima := 0.0
	scalingMaxima := 10.0

	parameters := slices.Concat(eden, d, sigma)
	parameters = append(parameters, delta, background, scaling)
	minimas := slices.Concat(edenMinima, dMinima, sigmaMinima)
	minimas = append(minimas, deltaMinima, backgroundMinima, scalingMinima)
	maximas := slices.Concat(edenMaxima, dMaxima, sigmaMaxima)
	maximas = append(maximas, deltaMaxima, backgroundMaxima, scalingMaxima)

	dataTracks := graphMap["intensity"].GetDataTracks()
	if len(dataTracks) == 0 {
		return "NO DATA POINTS IN GRAPH", nil
	}
	dataTrack := dataTracks[0]

	// Define error function
	penaltyFunction := func(params []float64) float64 {
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

		dataModel := dataTrack.Model(0, false)

		diff := 0.0
		for i := range dataModel {
			iy, err := intensityFunction.Eval(dataModel[i].X)
			if err != nil {
				fmt.Println("Error while calculating intensity:", err)
			}
			diff += math.Pow(dataModel[i].Y-iy, 2)
		}
		return diff
	}

	return "", minimizer.NewProblem(parameters, minimas, maximas, penaltyFunction, &minimizer.MinimiserConfig{
		LoopCount:     1e4,
		ParallelReads: true,
	})
}

// register functions which can be used for graph plotting
func registerFunctions() {
	functionMap["intensity"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
	functionMap["eden"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
	functionMap["test"] = function.NewFunction(function.Points{}, function.INTERPOLATION_NONE)
}

// creates the graph containers for the different graphs
func registerGraphs() *fyne.Container {
	graphMap["intensity"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Intensity Graph",
		IsLog:     true,
		AdaptDraw: true,
		Functions: function.Functions{functionMap["intensity"]},
		DisplayRange: &graph.GraphRange{
			Min: 0.01,
			Max: math.MaxFloat64,
		},
	})

	graphMap["eden"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Edensity Graph",
		IsLog:     false,
		AdaptDraw: false,
		Functions: function.Functions{functionMap["eden"]},
	})

	graphMap["test"] = graph.NewGraphCanvas(&graph.GraphConfig{
		Title:     "Test Graph",
		IsLog:     true,
		AdaptDraw: true,
		Functions: function.Functions{functionMap["test"]},
	})

	return container.NewGridWithColumns(2, graphMap["eden"], graphMap["intensity"], graphMap["test"])
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

	deltaQ, _ := param.Float("general", "deltaq", -0.000305927)
	background, _ := param.Float("general", "background", 1.43793e-7)
	scaling, _ := param.Float("general", "scaling", 0.888730)

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
		),   // top
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
// and set it to the intensity and eden graphs
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

	functionMap["intensity"].SetData(d)
	functionMap["eden"].SetData(d)
}

// RecalculateData recalculates the data for the intensity and eden graphs
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

	functionMap["intensity"].SetData(intensityPoints)
}
