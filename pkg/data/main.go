package data

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"math"
)

const (
	DefaultParameterName    string  = "Parameter"
	DefaultParameterValue   float64 = 10.0
	DefaultParameterDefault float64 = 10.0
	DefaultParameterMin     float64 = -math.MaxFloat64
	DefaultParameterMax     float64 = math.MaxFloat64
	DefaultParameterCheck   bool    = false
)

// from refl_monolayer.pro:780
var StartUpParameters = []*Parameter{
	NewTemplateParameter("Eden", "a", 0.0, DefaultParameterDefault),
	NewTemplateParameter("Eden", "1", 0.346197, DefaultParameterDefault),
	NewTemplateParameter("Eden", "2", 0.458849, DefaultParameterDefault),
	NewTemplateParameter("Eden", "b", 0.334000, DefaultParameterDefault),
	NewTemplateParameter("Roughness", "a/1", 3.39544, DefaultParameterDefault),
	NewTemplateParameter("Roughness", "1/2", 2.15980, DefaultParameterDefault),
	NewTemplateParameter("Roughness", "2/b", 3.90204, DefaultParameterDefault),
	NewTemplateParameter("Thickness", "1", 14.2657, DefaultParameterDefault),
	NewTemplateParameter("Thickness", "2", 10.6906, DefaultParameterDefault),
}

func NewSetString(val string) binding.String {
	b := binding.NewString()
	_ = b.Set(val)
	return b
}
func NewSetFloat(val float64) binding.Float {
	b := binding.NewFloat()
	_ = b.Set(val)
	return b
}

func NewTemplateParameter(group string, id string, value float64, defaultVal float64) *Parameter {
	name := fmt.Sprintf("%s-%s", group, id)
	p := &Parameter{
		DataChannel:    NewListenerGroup(),
		BindingChannel: NewChangeListenerGroup[binding.DataItem](),
		name:           NewSetString(name),
		value:          NewSetFloat(value),
		defaultValue:   NewSetFloat(defaultVal),
		min:            NewSetFloat(-math.MaxFloat64),
		max:            NewSetFloat(math.MaxFloat64),
		fixed:          binding.NewBool(),
		locked:         binding.NewBool(),
		class:          ParameterClass(group),
		uidType:        ParameterStaticID,
		uid:            ParameterID(name),
	}
	p.init()
	return p
}
