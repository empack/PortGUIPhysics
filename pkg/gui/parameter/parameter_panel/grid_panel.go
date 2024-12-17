package parameter_panel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"physicsGUI/pkg/data"
	"physicsGUI/pkg/gui/parameter"
	"slices"
)

type ParameterGridRenderer struct {
	fyne.WidgetRenderer
	impl *ParameterGrid
}

func NewParameterGridRenderer(impl *ParameterGrid) *ParameterGridRenderer {
	p := &ParameterGridRenderer{impl: impl}
	p.Update()
	return p
}

func (r *ParameterGridRenderer) Update() {
	objects := make([]fyne.CanvasObject, len(r.impl.objects)+1)
	for i, _ := range r.impl.objects {
		objects[i] = container.NewHBox(r.impl.objects[i], r.impl.options[i])
	}
	objects[len(r.impl.objects)] = r.impl.btnPnl
	cnt := container.NewAdaptiveGrid(r.impl.rowcol, objects...)
	r.WidgetRenderer = widget.NewSimpleRenderer(container.NewHScroll(cnt))
	r.WidgetRenderer.Layout(r.impl.Size())
}

type ParameterGrid struct {
	widget.BaseWidget
	objects  []*parameter.Wrapper
	options  []fyne.CanvasObject
	btnPnl   fyne.CanvasObject
	rowcol   int
	renderer *ParameterGridRenderer
	handler  *data.ParameterHandler
}

func (p *ParameterGrid) DataChanged(old *data.Parameter, new *data.Parameter) {
	if old != nil {
		for _, o := range p.objects {
			if o.Parameter == old {
				p.Remove(o)
			}
		}
	}
	if new != nil {
		p.Add(parameter.NewWrapper(new))
	}
}

func (p *ParameterGrid) CreateRenderer() fyne.WidgetRenderer {
	p.renderer = NewParameterGridRenderer(p)
	return p.renderer
}
func (p *ParameterGrid) Resize(size fyne.Size) {
	var maxParamWidth float32 = 0.1 // not 0.0 to prevent div by zero exception
	for i, c := range p.objects {
		minSize := c.MinSize()
		if minSize.Width+p.options[i].Size().Width > maxParamWidth {
			maxParamWidth = minSize.Width + p.options[i].Size().Width
		}
	}
	p.rowcol = int(size.Width / maxParamWidth)
	p.renderer.Update()
	p.BaseWidget.Resize(size)
}

func NewBoundParameterGrid(rowcol int, handler *data.ParameterHandler) *ParameterGrid {
	p := NewParameterGrid(rowcol)
	p.Bind(handler)
	return p
}

func NewParameterGrid(rowcol int, params ...*parameter.Wrapper) *ParameterGrid {
	objects := make([]fyne.CanvasObject, len(params))
	for i := 0; i < len(objects); i++ {
		objects[i] = params[i]
	}
	g := &ParameterGrid{
		BaseWidget: widget.BaseWidget{},
		objects:    params,
		rowcol:     rowcol,
	}
	g.btnPnl = container.NewStack(widget.NewButton("+", func() {
		g.Add(parameter.NewWrapper(data.NewParameter(data.ParameterDynamicID)))
	}))
	g.ExtendBaseWidget(g)
	return g
}

func (p *ParameterGrid) SetRowCols(rowcols int) {
	p.rowcol = rowcols
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterGrid) Add(parameter *parameter.Wrapper) {
	p.objects = append(p.objects, parameter)
	p.options = append(p.options, container.NewVBox(widget.NewButton("R", func() {
		if p.handler != nil {
			p.handler.Remove(parameter.Parameter)
		} else {
			p.Remove(parameter)
		}
	})))
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterGrid) Remove(parameter *parameter.Wrapper) {
	index := slices.Index(p.objects, parameter)
	if index != -1 {
		p.objects = append(p.objects[:index], p.objects[index+1:]...)
		p.options = append(p.options[:index], p.options[index+1:]...)
	}
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterGrid) Bind(handler *data.ParameterHandler) {
	p.handler = handler
	p.btnPnl = container.NewStack(widget.NewButton("+", func() {
		p.handler.Add(data.NewParameter(data.ParameterDynamicID))
	}))
	for i, _ := range p.options {
		p.options[i] = container.NewVBox(widget.NewButton("R", func() {
			p.handler.Remove(p.objects[i].Parameter)
		}))
	}
	handler.AddListener(p)
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterGrid) UnBind(handler *data.ParameterHandler) {
	p.handler = nil
	p.btnPnl = container.NewStack(widget.NewButton("+", func() {
		p.Add(parameter.NewWrapper(data.NewParameter(data.ParameterDynamicID)))
	}))
	for i, _ := range p.options {
		p.options[i] = container.NewVBox(widget.NewButton("R", func() {
			p.Remove(p.objects[i])
		}))
	}
	handler.RemoveListener(p)
}
