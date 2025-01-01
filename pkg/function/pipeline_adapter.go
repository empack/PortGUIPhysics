package function

import (
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/data/transformation"
)

type SldFunctionAdapter struct {
	transformation.BaseSegment[*data.ParameterHandler, Function]
}

func NewSldAdapter() *SldFunctionAdapter {
	return &SldFunctionAdapter{}
}

func (f *SldFunctionAdapter) Start(i int) {
	edensitys := f.In[i].GetByClass("Eden")
	d := f.In[i].GetByClass("Thickness")
	sigma := f.In[i].GetByClass("Roughness")
	f.Out[i] = NewSLDFunction(edensitys, d, sigma, 150).Function
	f.Wg.Done()
}
