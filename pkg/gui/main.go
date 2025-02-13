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
	"fyne.io/fyne/v2/widget"
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

		data := graphMap["intensity"].GetDataTracks()
		if len(data) != 1 {
			dialog.ShowError(errors.New("exactly 1 dataset in intensity graph needed"), MainWindow)
			return
		}

		if err := minimize(data[0], e1, e2, e3, e4, t1, t2, r1, r2, r3, delta, background, scaling); err != nil {
			fmt.Println("Error while minimizing:", err)
			dialog.ShowError(err, MainWindow)
			return
		}

		dialog.ShowInformation("Minimization Completed", fmt.Sprintf("Minimization TODO"), MainWindow)
		//dialog.ShowInformation("Minimization Completed", fmt.Sprintf("Minimization Stats:\n Error function calls: %f \n Remaining error: %f", minimum.Nfcn(), minimum.Fval()), MainWindow)
	})
}

func minimize(data *function.Function, parameters ...*param.Parameter[float64]) error {
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
	mFunc := minimizer.NewMinuitFcn(data, penaltyFunction, parameters)

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

func penaltyFunction(fcn *minimizer.MinuitFunction, params []float64) float64 {
	// parameter needed for parsing the parameters params[11] -> 12 parameters needed etc.
	paramCount := 12
	if len(params) != paramCount {
		dialog.ShowError(fmt.Errorf("penaltyFunction has %d parameters but expects %d", len(params), paramCount), MainWindow)
		return math.MaxFloat64
	}

	edenErr := params[0:4]
	dErr := params[4:6]
	sigmaErr := params[6:9]
	deltaErr := params[9]
	backgroundErr := params[10]
	scalingErr := params[11]

	log.Println("params", params)

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

	dataModel := fcn.Datatrack.Model(0, false)
	diff := 0.0
	for i := range dataModel {
		iy, err := intensityFunction.Eval(dataModel[i].X)
		if err != nil {
			fmt.Println("Error while calculating intensity:", err)
		}
		diff += math.Pow((dataModel[i].Y-iy)*math.Pow(dataModel[i].X*100, 4), 2)
	}
	return diff
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

	return container.NewGridWithColumns(2, graphMap["eden"], graphMap["intensity"])
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
