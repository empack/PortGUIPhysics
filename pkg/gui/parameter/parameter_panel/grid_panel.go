package parameter_panel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	cnt := container.NewAdaptiveGrid(r.impl.rowcol, r.impl.objects...)
	r.WidgetRenderer = widget.NewSimpleRenderer(container.NewHScroll(cnt))
	r.WidgetRenderer.Layout(r.impl.Size())
}

type ParameterGrid struct {
	widget.BaseWidget
	objects  []fyne.CanvasObject
	rowcol   int
	renderer *ParameterGridRenderer
}

func (p *ParameterGrid) CreateRenderer() fyne.WidgetRenderer {
	p.renderer = NewParameterGridRenderer(p)
	return p.renderer
}
func (p *ParameterGrid) Resize(size fyne.Size) {
	var maxParamWidth float32 = 0.1 // not 0.0 to prevent div by zero exception
	for _, c := range p.objects {
		minSize := c.MinSize()
		if minSize.Width > maxParamWidth {
			maxParamWidth = minSize.Width
		}
	}
	p.rowcol = int(size.Width / maxParamWidth)
	p.renderer.Update()
	p.BaseWidget.Resize(size)
}

func NewParameterGrid(rowcol int, parameter ...*parameter.Wrapper) *ParameterGrid {
	objects := make([]fyne.CanvasObject, len(parameter))
	for i := 0; i < len(objects); i++ {
		objects[i] = parameter[i]
	}
	g := &ParameterGrid{
		BaseWidget: widget.BaseWidget{},
		objects:    objects,
		rowcol:     rowcol,
	}
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
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterGrid) Remove(parameter *parameter.Wrapper) {
	index := slices.Index(p.objects, fyne.CanvasObject(parameter))
	if index != -1 {
		p.objects = append(p.objects[:index], p.objects[index+1:]...)
	}
	if p.renderer != nil {
		p.renderer.Update()
	}
}
