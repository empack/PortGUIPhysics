package function

import (
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/data/transformation"
)

type sldFunctionAdapter struct {
	transformation.BaseSegment[*data.ParameterHandler, Function]
}

func (f *SldFunction) Start(i int) {
	edensitys := f.In[i].GetByClass("Eden")
	d := f.In[i].GetByClass("Thickness")
	sigma := f.In[i].GetByClass("Roughness")
	f.Out = &NewSLDFunction(edensitys, d, sigma, 150).Function
}
