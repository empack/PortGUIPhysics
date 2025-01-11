package gui

import "physicsGUI/pkg/gui/graph"

var Graphs = []*graph.GraphCanvas{
	graph.NewGraphCanvas(&graph.GraphConfig{ // 0
		Title:      "Electron Density",
		IsLog:      false,
		Resolution: 150,
		Functions:  nil,
	}),
	graph.NewGraphCanvas(&graph.GraphConfig{ // 1
		Title:      "Intensity",
		IsLog:      false,
		Resolution: 150,
		Functions:  nil,
	}),
}
