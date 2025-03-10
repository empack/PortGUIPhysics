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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/davecgh/go-spew/spew"
	minuit "github.com/empack/minuit2go/pkg"
)

var (
	// App reference
	App        fyne.App
	MainWindow fyne.Window

	functionMap = make(map[string]*function.Function)
	graphMap    = make(map[string]*graph.GraphCanvas)
)

// adaption should not be necessary here
// Start GUI (function is blocking)
func Start() {
	App = app.NewWithID("GUI-Physics")
	App.Settings().SetTheme(theme.DarkTheme())
	MainWindow = App.NewWindow("Physics GUI")

	mainWindow()
}

// adaption should not be necessary
// onDrop is called when a file is dropped into the window
// imports the data if a file is dropped on a graph canvas
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
					newFunction := function.NewFunction(points)
					graphMap[mapIdentifier].AddDataTrack(newFunction)
					physics.AlterQZAxis(graphMap[mapIdentifier].GetDataTracks(), mapIdentifier)
				}
			}
			return
		}
	}
}

// adaption should not be necessary here
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

// adaption should not be necessary here
func createFileMenu() *fyne.Menu {
	mnLoad := fyne.NewMenuItem("Load", loadFileChooser)
	mnSave := fyne.NewMenuItem("Save", saveFileChooser)
	return fyne.NewMenu("File", mnLoad, mnSave)
}

// adaption should not be necessary here
// fit parameters to the experimental data passing everything to minuit
func minimize(experimentalData []function.Functions, parameters ...*param.Parameter[float64]) error {
	mnParams := minuit.NewEmptyMnUserParameters()

	var freeToChangeCnt int = 0

	for i, p := range parameters {
		if p == nil {
			return fmt.Errorf("minimizer: parameter %d is nil", i)
		}

		par, err := p.Get()
		if err != nil {
			return err
		}

		id := fmt.Sprintf("p%d", i)

		// if not checked, add as constant parameter
		if !p.IsChecked() {
			mnParams.Add(id, par)
			continue
		}

		min := p.GetRelative("min")
		max := p.GetRelative("max")

		// if min or max is nil, add as free parameter
		if min == nil || max == nil {
			mnParams.AddFree(id, par, 0.1)
			freeToChangeCnt++
			continue
		}

		// if min and max are set, add as limited parameter
		minV, err := min.Get()
		if err != nil {
			return err
		}

		maxV, err := max.Get()
		if err != nil {
			return err
		}

		mnParams.AddLimited(id, par, 0.1, minV, maxV)
		freeToChangeCnt++
	}

	if freeToChangeCnt == 0 {
		return fmt.Errorf("minimizer: Parameter to change selected")
	}

	// create minuit setup
	mFunc := minimizer.NewMinuitFcn(experimentalData, penaltyFunction, parameters)

	// create migrad
	migrad := minuit.NewMnMigradWithParameters(mFunc, mnParams)

	min, err := migrad.Minimize()
	if err != nil {
		return err
	}

	if !min.IsValid() {
		migrad2 := minuit.NewMnMigradWithParameterStateStrategy(mFunc, min.UserState(), minuit.NewMnStrategyWithStra(minuit.PreciseStrategy))
		min, err = migrad2.Minimize()
		if err != nil {
			return err
		}
	}

	fmt.Println("result")
	spew.Dump(min.UserParameters().Params())
	fmt.Printf("Fval: %f\n", min.Fval())
	fmt.Printf("FCNCall: %d\n", min.Nfcn())

	return mFunc.UpdateParameters(min.UserParameters().Params())
}

// adaption should not be necessary here
// mainWindow builds and renders the main GUI content, it will show and run the main window
func mainWindow() {
	registerFunctions()

	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				NewMinimizerControlPanel().Widget(),
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

	MainWindow.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("Program"),
		createFileMenu(),
	))
	MainWindow.Resize(fyne.NewSize(1000, 500))
	MainWindow.SetContent(content)
	MainWindow.SetOnDropped(onDrop)

	MainWindow.ShowAndRun()
}

//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! adapt everything from here !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// this is also the place where you need to pass:
// all current parameters and all experimental data tracks
func (this *MinimizerControlPanel) minimizerProblemSetup() error {
	// get parameters + experimental data and put them into minimize()
	edens := param.GetFloatGroup("eden")

	e1 := edens.GetParam("Eden a")
	e2 := edens.GetParam("Eden 1")
	e3 := edens.GetParam("Eden 2")
	e4 := edens.GetParam("Eden b")

	// get roughness parameters
	roughness := param.GetFloatGroup("rough")

	r1 := roughness.GetParam("Roughness a/1")
	r2 := roughness.GetParam("Roughness 1/2")
	r3 := roughness.GetParam("Roughness 2/b")

	// get thickness parameters
	thickness := param.GetFloatGroup("thick")

	t1 := thickness.GetParam("Thickness 1")
	t2 := thickness.GetParam("Thickness 2")

	// get general parameters
	general := param.GetFloatGroup("general")

	delta := general.GetParam("deltaq")
	background := general.GetParam("background")
	scaling := general.GetParam("scaling")

	experimentalData := []function.Functions{
		graphMap["intensity"].GetDataTracks(),
	}

	if err := this.minimize(experimentalData, e1, e2, e3, e4, t1, t2, r1, r2, r3, delta, background, scaling); err != nil {
		fmt.Println("Error while minimizing:", err)
		return err
	}
	return nil
}

// the penalty function defines the error we minimize with minuit
// !the order of the parameters needs to fit
func penaltyFunction(fcn *minimizer.MinuitFunction, params []float64) float64 {
	// parameter needed for parsing the parameters params[11] -> 12 parameters needed etc.
	paramCount := 12
	if len(params) != paramCount {
		dialog.ShowError(fmt.Errorf("penaltyFunction has %d parameters but expects %d", len(params), paramCount), MainWindow)
		return math.MaxFloat64
	}

	//sort the parameters
	edenErr := params[0:4]
	dErr := params[4:6]
	sigmaErr := params[6:9]
	deltaErr := params[9]
	backgroundErr := params[10]
	scalingErr := params[11]

	log.Println("params", params)

	//precalculation for intensities
	edenPoints, err := physics.GetEdensities(edenErr, dErr, sigmaErr)
	if err != nil {
		fmt.Println("Error while calculating edensities:", err)
		return math.MaxFloat64
	}

	//intensity calculation itself
	intensityPoints := physics.CalculateIntensityPoints(edenPoints, deltaErr, &physics.IntensityOptions{
		Background: backgroundErr,
		Scaling:    scalingErr,
	})
	dataPoints := make([]function.Points, len(fcn.ExperimentalData))
	for _, funcInGraph := range fcn.ExperimentalData {
		for i, expData := range funcInGraph {
			dataPoints[i] = expData.GetData()
		}
	}

	//real penalty calculation
	diff, err := physics.Sim2SigRMS(dataPoints, intensityPoints)
	if err != nil {
		dialog.ShowError(err, MainWindow)
	}

	return diff
}

// register functions which can be used for graph plotting
// this is the place to add function plots which are shown in graphs
func registerFunctions() {
	//a function needs to be added to the functionMap using a unique identifier so we can further handle it
	//interpolation mode can rather be ignored here
	functionMap["intensity"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
	functionMap["eden"] = function.NewEmptyFunction(function.INTERPOLATION_NONE)
}

// creates the graph containers for the different graphs
// this is the place to add graphs to the GUI
func registerGraphs() *fyne.Container {

	//a graph needs to be added to the graphMap using a unique identifier so we can further handle it
	graphMap["intensity"] = graph.NewGraphCanvas(&graph.GraphConfig{

		//title shown inside the GUI
		Title: "Intensity Graph",

		//use logarithmic scaling (both x and y axis)
		IsLog: true,

		//use magic scaling: p.Y = math.Pow(p.X, 4) * p.Y;  p.Error = math.Pow(p.X, 4) * p.Error
		AdaptDraw: true,

		//chose the function to show inside the graph by it's identifier
		Functions: function.Functions{functionMap["intensity"]},

		//optionally set an x-range to plot, points outside it are ignored
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

	//chose how you like to arrange the graphs insige the GUI
	return container.NewGridWithColumns(2, graphMap["eden"], graphMap["intensity"])
}

// creates and registers the parameter and adds them to the parameter repository
// this is the place to alter parameters:
func registerParams() *fyne.Container {
	//created with a group name, an individual name and a default value
	//you can get parameters as a group or individually (combining group and individual name) later on
	//this can be helpful to easily pass similar parameters to a function and iterate over them
	edenA, _ := param.FloatMinMax("eden", "Eden a", 0.0)
	eden1, _ := param.FloatMinMax("eden", "Eden 1", 0.346197)
	eden2, _ := param.FloatMinMax("eden", "Eden 2", 0.458849)
	edenB, _ := param.FloatMinMax("eden", "Eden b", 0.334000)

	roughnessA1, _ := param.FloatMinMax("rough", "Roughness a/1", 3.39544)
	roughness12, _ := param.FloatMinMax("rough", "Roughness 1/2", 2.15980)
	roughness2B, _ := param.FloatMinMax("rough", "Roughness 2/b", 3.90204)

	thickness1, _ := param.FloatMinMax("thick", "Thickness 1", 14.2657)
	thickness2, _ := param.FloatMinMax("thick", "Thickness 2", 10.6906)

	//parameters can be created with (above) or without (below) two additional fields for minimum and maximum values
	deltaQ, _ := param.Float("general", "deltaq", -0.000305927)
	background, _ := param.Float("general", "background", 1.43793e-7)
	scaling, _ := param.Float("general", "scaling", 0.888730)

	//you can chose how to arrange the parameters inside the GUI here
	//by now it's 4x4 partitioning where some partitions are left empty
	containers := container.NewVBox(
		container.NewGridWithColumns(4, edenA, eden1, eden2, edenB),
		container.NewGridWithColumns(4, roughnessA1, roughness12, roughness2B),
		container.NewGridWithColumns(4, thickness1, thickness2),
		container.NewGridWithColumns(4, deltaQ, background, scaling),
	)

	//makes a scrollbar for the parameters
	con2 := container.NewScroll(containers)
	con2.SetMinSize(fyne.NewSize(300, 300))
	return container.NewStack(con2)
}

// Insert your adapted physical calculations and parameters here!
// RecalculateData recalculates the data for the current graphs
// current parameter values need to be fetched, the physical calculations done and resulting points set to the functions
func RecalculateData() {
	// Fetch all parameters here
	// Get current parameters by group identifier
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

	// get general parameters individually
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

	// calculate all functions which need to be updated here

	// calculate edensities
	edenPoints, err := physics.GetEdensities(eden, d, sigma)
	//only potential error handling
	if err != nil {
		log.Println("Error while calculating edensities:", err)
		return
	} else {
		//set points to function which is automatically shown inside the graph
		functionMap["eden"].SetData(edenPoints)
	}

	// calculate intensities
	intensityPoints := physics.CalculateIntensityPoints(edenPoints, delta, &physics.IntensityOptions{
		Background: background,
		Scaling:    scaling,
	})
	//set points to function which is automatically shown inside the graph
	functionMap["intensity"].SetData(intensityPoints)
}
