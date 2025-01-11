package gui

import (
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/function"
	"physicsGUI/pkg/gui/graph"
)

var ParameterList = data.ParameterList
var Helper = data.Helper

var Graphs = []*graph.GraphCanvas{
	graph.NewGraphCanvas(&graph.GraphConfig{ // 0
		Title:      "Electron Density",
		IsLog:      false,
		Resolution: 150,
		Functions:  nil,
		Calculation: []func() function.Points{
			func() function.Points {
				return Helper.GetEden(
					[]*data.Parameter{&ParameterList[0], &ParameterList[1], &ParameterList[2], &ParameterList[3]},
					[]*data.Parameter{&ParameterList[7], &ParameterList[8]},
					[]*data.Parameter{&ParameterList[4], &ParameterList[5], &ParameterList[6]},
					data.ZNumber)
			},
		},
	}),
	graph.NewGraphCanvas(&graph.GraphConfig{ // 1
		Title:      "Intensity",
		IsLog:      true,
		Resolution: 150,
		Functions:  nil,
		Calculation: []func() function.Points{
			func() function.Points {
				return Helper.GetIntensity(Helper.Get_zaxis(Helper.GetValues(&ParameterList[7], &ParameterList[8]), data.ZNumber), Helper.GetEden(
					[]*data.Parameter{&ParameterList[0], &ParameterList[1], &ParameterList[2], &ParameterList[3]},
					[]*data.Parameter{&ParameterList[7], &ParameterList[8]},
					[]*data.Parameter{&ParameterList[4], &ParameterList[5], &ParameterList[6]},
					data.ZNumber), &ParameterList[9], &ParameterList[10], &ParameterList[11])
			},
		},
	}),
}
