package parameter

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/gui/parameter/custom_bindings"
)

type FieldListener struct {
	field    binding.DataItem
	listener binding.DataListener
}

type Wrapper struct {
	widget.BaseWidget
	check         *widget.Check
	name          *widget.Entry
	val           *FilteredEntry
	min           *FilteredEntry
	max           *FilteredEntry
	parameter     *data.Parameter
	fieldListener []*FieldListener
}

func (p *Wrapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewVBox(p.name, container.NewHBox(container.NewCenter(p.check), container.NewCenter(p.val), container.NewVBox(p.max, p.min))))
}
func (p *Wrapper) MinSize() fyne.Size {
	var minWidth, minHeight float32 = 20.0, 0.0 // padding offset
	minWidth += p.check.Size().Width
	minWidth += p.val.MinSize().Width
	minWidth += max(p.min.MinSize().Width, p.max.MinSize().Width)

	minHeight += p.name.MinSize().Height
	minHeight += max(p.check.Size().Height, p.val.MinSize().Height, p.min.MinSize().Height+p.max.MinSize().Height)
	return fyne.Size{
		Width:  max(minWidth, p.name.MinSize().Width),
		Height: minHeight,
	}
}

func NewWrapper(parameter *data.Parameter) *Wrapper {
	// create name text field with linked data
	name := widget.NewEntry()
	name.Validator = nil
	// create filtered entry fields, which only accept runes relevant for float inputs
	val := NewFilteredEntry('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', 'e', '.')
	minV := NewFilteredEntry('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', 'e', '.')
	minV.PlaceHolder = "Min"
	maxV := NewFilteredEntry('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', 'e', '.')
	maxV.PlaceHolder = "Max"

	// make checkbox for locking data for minimizer
	check := widget.NewCheck("", func(b bool) {
		//TODO add icon and icon change on pressing?
	})
	p := &Wrapper{
		BaseWidget: widget.BaseWidget{},
		name:       name,
		check:      check,
		val:        val,
		min:        minV,
		max:        maxV,
		parameter:  parameter,
	}

	// define Listeners
	defaultUpdater := binding.NewDataListener(func() {
		// set placeholder text of val to default value, when available
		if def, err := parameter.GetDefault().Get(); err == nil {
			val.PlaceHolder = fmt.Sprint(def)
		}
	})
	p.addInputListener(parameter.GetDefault(), defaultUpdater) // register listener with relevant binding for update on binding change
	disableInputFieldsListener := binding.NewDataListener(func() {
		// update disable input fields, when parameter is set to locked/unlocked
		if locked, err := parameter.GetLocked().Get(); err == nil {
			if locked {
				p.name.Disable()
				p.check.Disable()
				p.val.Disable()
				p.min.Disable()
				p.max.Disable()
			} else {
				p.name.Enable()
				p.check.Enable()
				p.val.Enable()
				p.min.Enable()
				p.max.Enable()
			}
		}
	})
	p.addInputListener(parameter.GetLocked(), disableInputFieldsListener) // register listener with relevant binding for update on binding change

	// Bind gui representation to the current data in Parameter
	p.rebind()

	// set to default value, if value is submitted empty
	p.val.OnSubmitted = func(s string) {
		if s == "" {
			if get, err := p.parameter.GetDefault().Get(); err == nil {
				p.val.SetText(fmt.Sprint(get))
			}
		}
	}

	// Move binding from old binding to new binding, when binding changed and update data bindings
	parameter.BindingChannel.AddListener(data.NewChangeListener(func(old binding.DataItem, new binding.DataItem) {
		for i, l := range p.fieldListener {
			if l.field == old {
				old.RemoveListener(l.listener)
				new.AddListener(l.listener)
				p.fieldListener[i].field = new
			}
		}
		p.rebind()
	}))

	p.ExtendBaseWidget(p)
	return p
}

func (p *Wrapper) addInputListener(field binding.DataItem, listenerF binding.DataListener) {
	p.fieldListener = append(p.fieldListener, &FieldListener{
		field:    field,
		listener: listenerF,
	})
	field.AddListener(listenerF)

}

func (p *Wrapper) rebind() {
	p.name.Bind(p.parameter.GetName())
	p.val.Bind(custom_bindings.NewLazyFloatToString(p.parameter.GetValue(), p.parameter.GetDefault()))
	p.min.Bind(custom_bindings.NewLazyFloatToString(p.parameter.GetMin(), nil))
	p.max.Bind(custom_bindings.NewLazyFloatToString(p.parameter.GetMax(), nil))
	p.check.Bind(p.parameter.GetFixed())
}
