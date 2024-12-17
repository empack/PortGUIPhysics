package data

import (
	"fyne.io/fyne/v2/data/binding"
)

type ParameterIDType int
type ParameterClass binding.String
type ParameterID string

const (
	ParameterStaticID  = ParameterIDType(iota)
	ParameterDynamicID = ParameterIDType(iota)
)

type Parameter struct {
	DataChannel    *ListenerGroup                         // gets notified, when data in bindings changed
	BindingChannel *ChangeListenerGroup[binding.DataItem] // gets notified, when bindings changed
	name           binding.String
	value          binding.Float
	defaultValue   binding.Float
	min            binding.Float
	max            binding.Float
	fixed          binding.Bool // value should not be changed by minimizer
	locked         binding.Bool // disable input fields
	class          ParameterClass
	uidType        ParameterIDType
	uid            ParameterID
}

func NewParameter(uidType ParameterIDType) *Parameter {
	p := &Parameter{
		DataChannel:    NewListenerGroup(),
		BindingChannel: NewChangeListenerGroup[binding.DataItem](),
		uidType:        uidType,
		uid:            ParameterID("Parameter"),
	}
	name := binding.NewString()
	_ = name.Set(DefaultParameterName)
	p.BindName(name)
	value := binding.NewFloat()
	_ = value.Set(DefaultParameterValue)
	p.BindValue(value)
	minV := binding.NewFloat()
	_ = minV.Set(DefaultParameterMin)
	p.BindMin(minV)
	maxV := binding.NewFloat()
	_ = maxV.Set(DefaultParameterMax)
	p.BindMax(maxV)
	defaultValue := binding.NewFloat()
	_ = defaultValue.Set(DefaultParameterDefault)
	p.BindDefault(defaultValue)
	fixed := binding.NewBool()
	_ = fixed.Set(DefaultParameterCheck)
	p.BindFixed(fixed)
	p.BindLocked(binding.NewBool())

	return p
}

// Get Binding Functions

func (p *Parameter) GetName() binding.String {
	return p.name
}
func (p *Parameter) GetValue() binding.Float {
	return p.value
}
func (p *Parameter) GetMin() binding.Float {
	return p.min
}
func (p *Parameter) GetMax() binding.Float {
	return p.max
}
func (p *Parameter) GetDefault() binding.Float {
	return p.defaultValue
}
func (p *Parameter) GetLocked() binding.Bool {
	return p.locked
}
func (p *Parameter) GetFixed() binding.Bool {
	return p.fixed
}

// Set Binding Functions

func (p *Parameter) BindName(newBinding binding.String) {
	if p.name != nil {
		p.name.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.name, newBinding)
	p.name = newBinding
	p.name.AddListener(p.DataChannel)
}
func (p *Parameter) BindValue(newBinding binding.Float) {
	if p.value != nil {
		p.value.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.value, newBinding)
	p.value = newBinding
	p.value.AddListener(p.DataChannel)
}
func (p *Parameter) BindMin(newBinding binding.Float) {
	if p.min != nil {
		p.min.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.min, newBinding)
	p.min = newBinding
	p.min.AddListener(p.DataChannel)
}
func (p *Parameter) BindMax(newBinding binding.Float) {
	if p.max != nil {
		p.max.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.max, newBinding)
	p.max = newBinding
	p.max.AddListener(p.DataChannel)
}
func (p *Parameter) BindDefault(newBinding binding.Float) {
	if p.defaultValue != nil {
		p.defaultValue.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.defaultValue, newBinding)
	p.defaultValue = newBinding
	p.defaultValue.AddListener(p.DataChannel)
}
func (p *Parameter) BindFixed(newBinding binding.Bool) {
	if p.fixed != nil {
		p.fixed.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.fixed, newBinding)
	p.fixed = newBinding
	p.fixed.AddListener(p.DataChannel)
}
func (p *Parameter) BindLocked(newBinding binding.Bool) {
	if p.locked != nil {
		p.locked.RemoveListener(p.DataChannel)
	}
	p.BindingChannel.trigger(p.locked, newBinding)
	p.locked = newBinding
	p.locked.AddListener(p.DataChannel)
}
