package graph

import (
	"physicsGUI/pkg/data/transformation"
	"physicsGUI/pkg/function"
)

type pipelineCanvasAdapter struct {
	transformation.BaseSegment[function.Function, function.Function]
}

func (g *GraphCanvas) Start(i int) {
	f := g.In[i]
	g.UpdateFunction(&f)
	g.Out = &f
	g.Wg.Done()
}
