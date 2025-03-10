package gui

import (
	"fmt"
	"math"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/helper"
	"physicsGUI/pkg/gui/param"
	"physicsGUI/pkg/minimizer"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	minuit "github.com/empack/minuit2go/pkg"
)

type MinimizerState int

const (
	MinimizerNotStarted = MinimizerState(0x1 << iota)
	MinimizerRunning    = MinimizerState(0x1 << iota)
	MinimizerPaused     = MinimizerState(0x1 << iota)
	MinimizerFinished   = MinimizerState(0x1 << iota)
	MinimizerFailed     = MinimizerState(0x1 << iota)
)

type SharedMinimizerData struct {
	rw          sync.RWMutex
	threadCount int
	mnParams    *minuit.MnUserParameters
	mFunc       *minimizer.MinuitFunction
	err         error
}

func (this *MinimizerControlPanel) MinuitUpdateHandler() {
	var migrad *minuit.MnMigrad
	var migrad2 *minuit.MnMigrad
	lastError := math.MaxFloat64

	stateReader := func() MinimizerState {
		this.rw.RLock()
		defer this.rw.RUnlock()
		return this.state
	}

	for {
		switch stateReader() {
		case MinimizerRunning:
			this.sharedStorage.rw.Lock()
			if migrad == nil {
				// create migrad
				migrad = minuit.NewMnMigradWithParameters(this.sharedStorage.mFunc, this.sharedStorage.mnParams)
			}

			res, err := migrad.MinimizeWithMaxfcn(50)

			if err != nil {
				this.SetStats(err, 0, 0)
				this.sharedStorage.err = err
				migrad = nil
				migrad2 = nil
				this.sharedStorage.rw.Unlock()
				this.Failed(err)
				continue
			}

			if !res.IsValid() {
				if migrad2 == nil {
					migrad2 = minuit.NewMnMigradWithParameterStateStrategy(this.sharedStorage.mFunc, res.UserState(), minuit.NewMnStrategyWithStra(minuit.PreciseStrategy))
				}
				res, err = migrad2.MinimizeWithMaxfcn(50)
				if err != nil {
					this.SetStats(err, 0, 0)
					this.sharedStorage.err = err
					migrad = nil
					migrad2 = nil
					this.sharedStorage.rw.Unlock()
					this.Failed(err)
					continue
				}
			}
			this.SetStats(err, res.Fval(), res.Nfcn())
			_ = this.sharedStorage.mFunc.UpdateParameters(res.UserParameters().Params())
			this.sharedStorage.rw.Unlock()
			if res.Fval() == lastError {
				migrad = nil
				migrad2 = nil
				lastError = math.MaxFloat64
				this.Completed()
			} else {
				lastError = res.Fval()
			}
			continue
		case MinimizerPaused:
			time.Sleep(500 * time.Millisecond) // TODO idle efficient and responsive
			continue
		case MinimizerFinished | MinimizerFailed | MinimizerNotStarted:
			migrad = nil
			migrad2 = nil
			lastError = math.MaxFloat64
			time.Sleep(1000 * time.Millisecond) // TODO idle efficient and responsive
			continue
		}
	}
}

type MinimizerControlPanel struct {
	rw               sync.RWMutex
	state            MinimizerState
	btnPause         *widget.Button
	btnContinue      *widget.Button
	btnStop          *widget.Button
	btnStart         *widget.Button
	lblNCalls        *widget.Label
	lblFVal          *widget.Label
	lblError         *widget.Label
	lblStatus        *widget.Label
	oldMinimizerData []float64
	sharedStorage    *SharedMinimizerData
}

func NewMinimizerControlPanel() *MinimizerControlPanel {
	pnlControl := &MinimizerControlPanel{
		rw:               sync.RWMutex{},
		state:            MinimizerNotStarted,
		btnPause:         nil,
		btnContinue:      nil,
		btnStop:          nil,
		btnStart:         nil,
		lblNCalls:        widget.NewLabel("Calls: -"),
		lblFVal:          widget.NewLabel("FVal: -"),
		lblError:         widget.NewLabel(""),
		lblStatus:        widget.NewLabel("Not Initialized"),
		oldMinimizerData: nil,
		sharedStorage:    &SharedMinimizerData{},
	}
	go pnlControl.MinuitUpdateHandler()
	pnlControl.lblError.Hide()
	pnlControl.btnPause = widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), pnlControl.Pause)
	pnlControl.btnPause.Disable()
	pnlControl.btnPause.Hide()
	pnlControl.btnContinue = widget.NewButtonWithIcon("Resume", theme.NavigateNextIcon(), pnlControl.Continue)
	pnlControl.btnContinue.Disable()
	pnlControl.btnContinue.Hide()
	pnlControl.btnStart = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), pnlControl.Start)
	pnlControl.btnStop = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), pnlControl.Stop)
	pnlControl.btnStop.Disable()
	pnlControl.btnStop.Hide()
	return pnlControl
}

func (this *MinimizerControlPanel) Widget() fyne.CanvasObject {
	return container.NewHBox(this.btnStart, this.btnContinue, this.btnPause, this.btnStop, helper.CreateSeparator(), container.NewVBox(container.NewHBox(this.lblError, this.lblFVal, this.lblNCalls), helper.CreateSeparator(), this.lblStatus))
}

func (this *MinimizerControlPanel) Pause() {
	if this.state == MinimizerRunning {
		this.btnStart.Disable()
		this.btnStart.Hide()
		this.btnPause.Disable()
		this.btnPause.Hide()
		this.btnContinue.Disable()
		this.btnContinue.Hide()
		this.btnStop.Disable()
		this.btnStop.Hide()

		this.state = MinimizerPaused

		// this blocks until current cycle is completed
		this.sharedStorage.rw.Lock()
		this.sharedStorage.rw.Unlock()

		this.lblStatus.SetText("Paused")

		this.btnContinue.Enable()
		this.btnContinue.Show()
		this.btnStop.Enable()
		this.btnStop.Show()

	}
}

func (this *MinimizerControlPanel) Continue() {
	if this.state == MinimizerPaused {
		this.btnStart.Disable()
		this.btnStart.Hide()
		this.btnPause.Disable()
		this.btnPause.Hide()
		this.btnContinue.Disable()
		this.btnContinue.Hide()
		this.btnStop.Disable()
		this.btnStop.Hide()

		this.state = MinimizerRunning

		// this blocks until current cycle is completed
		this.sharedStorage.rw.Lock()
		this.sharedStorage.rw.Unlock()
		this.lblStatus.SetText("Running")

		this.btnPause.Enable()
		this.btnPause.Show()
		this.btnStop.Enable()
		this.btnStop.Show()
	}
}

func (this *MinimizerControlPanel) Start() {
	if this.state == MinimizerNotStarted || this.state == MinimizerFinished || this.state == MinimizerFailed {
		err := this.minimizerProblemSetup()
		if err != nil {
			dialog.ShowError(err, MainWindow)
			return
		}
		this.sharedStorage.rw.RLock()
		this.oldMinimizerData = this.sharedStorage.mnParams.Params()
		this.sharedStorage.rw.RUnlock()

		this.btnStart.Disable()
		this.btnStart.Hide()
		this.btnPause.Disable()
		this.btnPause.Hide()
		this.btnContinue.Disable()
		this.btnContinue.Hide()
		this.btnStop.Disable()
		this.btnStop.Hide()

		this.state = MinimizerPaused

		// this blocks until current cycle is completed
		this.sharedStorage.rw.Lock()
		this.sharedStorage.rw.Unlock()
		this.lblStatus.SetText("Ready")

		this.Continue()
	}
}

func (this *MinimizerControlPanel) Stop() {
	this.btnStart.Disable()
	this.btnStart.Hide()
	this.btnPause.Disable()
	this.btnPause.Hide()
	this.btnContinue.Disable()
	this.btnContinue.Hide()
	this.btnStop.Disable()
	this.btnStop.Hide()

	this.Reset()
	this.SetStats(nil, 0, 0)

	this.sharedStorage.rw.Lock()
	_ = this.sharedStorage.mFunc.UpdateParameters(this.oldMinimizerData)
	this.sharedStorage.rw.Unlock()
	this.lblStatus.SetText("Not Initialized")
}

func (this *MinimizerControlPanel) Reset() {
	this.btnStart.Disable()
	this.btnStart.Hide()
	this.btnPause.Disable()
	this.btnPause.Hide()
	this.btnContinue.Disable()
	this.btnContinue.Hide()
	this.btnStop.Disable()
	this.btnStop.Hide()

	this.state = MinimizerNotStarted

	// this blocks until current cycle is completed
	this.sharedStorage.rw.Lock()
	this.sharedStorage.rw.Unlock()
	this.lblStatus.SetText("Not Initialized")
	this.btnStart.Enable()
	this.btnStart.Show()

}

func (this *MinimizerControlPanel) Completed() {
	this.Reset()
	dialog.ShowInformation("Minimizer Completed", "Minimizer finished. No further improvements found.", MainWindow)
	this.state = MinimizerFinished
	// this blocks until current cycle is completed
	this.sharedStorage.rw.Lock()
	this.sharedStorage.rw.Unlock()
	this.lblStatus.SetText("Completed")
}

func (this *MinimizerControlPanel) Failed(err error) {
	this.Reset()
	dialog.ShowError(err, MainWindow)
	//TODO ask user if he wants to use data?
	this.state = MinimizerFailed
	// this blocks until current cycle is completed
	this.sharedStorage.rw.Lock()
	this.sharedStorage.rw.Unlock()
	this.lblStatus.SetText("Failed")
}

func (this *MinimizerControlPanel) minimize(experimentalData []function.Functions, parameters ...*param.Parameter[float64]) error {
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
		return fmt.Errorf("minimizer: No parameter(s) selected to be minimized")
	}

	// create minuit setup
	mFunc := minimizer.NewMinuitFcn(experimentalData, penaltyFunction, parameters)

	this.sharedStorage.rw.Lock()
	this.sharedStorage.mFunc = mFunc
	this.sharedStorage.mnParams = mnParams
	this.sharedStorage.rw.Unlock()

	return nil
}

func (this *MinimizerControlPanel) SetStats(err error, fVal float64, nCalls int) {
	if err != nil {
		this.lblError.SetText(err.Error())
		this.lblError.Show()
		this.lblFVal.Hide()
		this.lblNCalls.Hide()
	} else {
		this.lblError.SetText("")
		this.lblError.Hide()
		this.lblFVal.Show()
		this.lblNCalls.Show()
	}
	if fVal == 0 {
		this.lblFVal.SetText(fmt.Sprintf("FVal: -"))
	} else {
		this.lblFVal.SetText(fmt.Sprintf("FVal: %g", fVal))
	}
	if nCalls == 0 {
		this.lblNCalls.SetText(fmt.Sprintf("Calls: -"))
	} else {
		this.lblNCalls.SetText(fmt.Sprintf("Calls: %d", nCalls))
	}
}
