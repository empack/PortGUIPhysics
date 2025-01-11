package data

import (
	"physicsGUI/pkg/function"
	mod_intensity "physicsGUI/pkg/physics/intensity"
)

const r_e = 2.81e-5 // classical electron radius in angstrom
const zNumber = 150

type ifunctions struct {
	get_zaxis    func(d []float64, zNumber int) []float64
	getValues    func(p ...*Parameter) []float64
	getEden      func(eden []*Parameter, d []*Parameter, sigma []*Parameter, zNumber int) function.Points
	getIntensity func(qzaxis []float64, sld function.Points, deltaz *Parameter, background *Parameter, scaling *Parameter) function.Points
}

var functions = ifunctions{
	get_zaxis: get_zaxis,
	getValues: func(p ...*Parameter) []float64 {
		res := make([]float64, len(p))
		for i := range res {
			if v, err := p[i].Val.Get(); err == nil {
				res[i] = v
			} else {
				panic(err)
			}
		}
		return res
	},
	getEden: func(eden []*Parameter, d []*Parameter, sigma []*Parameter, zNumber int) function.Points {
		edenv := make([]float64, len(eden))
		dv := make([]float64, len(d))
		sigmav := make([]float64, len(sigma))
		for i := range eden {
			if v, err := eden[i].Val.Get(); err == nil {
				edenv[i] = v
			} else {
				panic(v)
			}
		}
		for i := range d {
			if v, err := d[i].Val.Get(); err == nil {
				dv[i] = v
			} else {
				panic(v)
			}
		}
		for i := range sigma {
			if v, err := sigma[i].Val.Get(); err == nil {
				sigmav[i] = v
			} else {
				panic(v)
			}
		}
		if points, err := getEden(edenv, dv, sigmav, zNumber); err == nil {
			return points
		}
		panic("Eden Wrapper crashed")
	},
	getIntensity: func(qzaxis []float64, sld function.Points, deltaz *Parameter, background *Parameter, scaling *Parameter) function.Points {
		deltazv, err := deltaz.Val.Get()
		if err != nil {
			panic(err)
		}
		sldv := make([]float64, len(sld))
		for i := 0; i < len(sld); i++ {
			sldv[i] = r_e * sld[i].Y
		}

		res := make(function.Points, len(sld))
		backgroundv, err := background.Val.Get()
		if err != nil {
			panic(err)
		}
		scalingv, err := scaling.Val.Get()
		if err != nil {
			panic(err)
		}
		opts := &mod_intensity.IntensityOptions{
			Background: backgroundv,
			Scaling:    scalingv,
		}
		intensityv := mod_intensity.CalculateIntensity(qzaxis, deltazv, sldv, opts)
		for i := range intensityv {
			res[i] = &function.Point{
				X:     sld[i].X,
				Y:     intensityv[i],
				Error: sld[i].Error,
			}
		}
		return res
	},
}

type PlotUpdate struct {
	plots []function.Points
}

func defineFunctions(plotUpdate *PlotUpdate) {
	eden := functions.getEden(
		[]*Parameter{&ParameterList[0], &ParameterList[1], &ParameterList[2], &ParameterList[3]},
		[]*Parameter{&ParameterList[7], &ParameterList[8]},
		[]*Parameter{&ParameterList[4], &ParameterList[5], &ParameterList[6]},
		zNumber)
	qzaxis := functions.get_zaxis(functions.getValues(&ParameterList[7], &ParameterList[8]), zNumber)
	intensity := functions.getIntensity(qzaxis, eden, &ParameterList[9], &ParameterList[10], &ParameterList[11])

	plotUpdate.plots[0] = eden
	plotUpdate.plots[1] = intensity
}
