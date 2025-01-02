package function

import (
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/data/transformation"
	"strings"
)

type SldFunctionAdapter struct {
	transformation.BaseSegment[[]*data.Parameter, Function]
}

func NewSldAdapter() *SldFunctionAdapter {
	return &SldFunctionAdapter{}
}

func (f *SldFunctionAdapter) Start(i int) {
	edensitys := make([]*data.Parameter, 0)
	sigma := make([]*data.Parameter, 0)
	d := make([]*data.Parameter, 0)
	for _, param := range f.In[0] {
		if strings.ToLower(string(param.GetClass())) == strings.ToLower("Eden") {
			edensitys = append(edensitys, param)
		} else if strings.ToLower(string(param.GetClass())) == strings.ToLower("Roughness") {
			sigma = append(sigma, param)
		} else if strings.ToLower(string(param.GetClass())) == strings.ToLower("Thickness") {
			d = append(d, param)
		}
	}
	f.Out[i] = NewSLDFunction(edensitys, d, sigma, 150).Function
	f.Wg.Done()
}
