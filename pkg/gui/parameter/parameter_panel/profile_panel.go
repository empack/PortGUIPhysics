package parameter_panel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"physicsGUI/pkg/gui/parameter"
	"slices"
)

type ParameterProfileRenderer struct {
	fyne.WidgetRenderer
	impl *ParameterProfile
}

func NewParameterProfileRenderer(impl *ParameterProfile) *ParameterProfileRenderer {
	p := &ParameterProfileRenderer{impl: impl}
	p.Update()
	return p
}

func (r *ParameterProfileRenderer) Update() {
	topPanel := container.NewBorder(nil, nil, nil, r.impl.buttonPanel, r.impl.name)
	var vObjectList = make([]fyne.CanvasObject, len(r.impl.parameter)+2)
	vObjectList[0] = topPanel
	vObjectList[1] = widget.NewSeparator()
	for i := 2; i < len(vObjectList); i++ {
		vObjectList[i] = r.impl.parameter[i-2]
	}
	contentPanel := container.NewVScroll(container.NewVBox(vObjectList...))
	r.WidgetRenderer = widget.NewSimpleRenderer(contentPanel)
	r.WidgetRenderer.Layout(r.impl.Size())
}

type ParameterProfile struct {
	widget.BaseWidget
	name        *widget.Entry
	buttonPanel fyne.CanvasObject
	parameter   []*parameter.ParameterWrapper
	renderer    *ParameterProfileRenderer
}

func (p *ParameterProfile) CreateRenderer() fyne.WidgetRenderer {
	p.renderer = NewParameterProfileRenderer(p)
	return p.renderer
}

func NewParameterProfile(nameVal binding.String) *ParameterProfile {
	name := widget.NewEntryWithData(nameVal)
	name.Validator = nil
	p := &ParameterProfile{
		name:      name,
		parameter: []*parameter.ParameterWrapper{},
	}
	p.ExtendBaseWidget(p)
	return p
}
func (p *ParameterProfile) Add(parameter *parameter.ParameterWrapper) {
	p.parameter = append(p.parameter, parameter)
	if p.renderer != nil {
		p.renderer.Update()
	}
}
func (p *ParameterProfile) Remove(parameter *parameter.ParameterWrapper) {
	parameterIndex := slices.Index(p.parameter, parameter)
	if parameterIndex == -1 {
		return
	}
	p.parameter = append(p.parameter[:parameterIndex], p.parameter[parameterIndex+1:]...)
	if p.renderer != nil {
		p.renderer.Update()
	}
}
func (p *ParameterProfile) SetButtonPanel(pnl fyne.CanvasObject) {
	p.buttonPanel = pnl
	if p.renderer != nil {
		p.renderer.Update()
	}
}

type ParameterProfilePanelRenderer struct {
	fyne.WidgetRenderer
	impl *ParameterProfilePanel
}

func NewParameterProfilePanelRenderer(impl *ParameterProfilePanel) *ParameterProfilePanelRenderer {
	r := &ParameterProfilePanelRenderer{impl: impl}
	r.Update()
	return r
}

func (r *ParameterProfilePanelRenderer) Update() {
	// prepare option panel, if options available
	options := container.NewStack()
	if r.impl.optSettings != nil && len(r.impl.optSettings) > 0 {
		opts := make([]fyne.CanvasObject, len(r.impl.optSettings))
		for i := 0; i < len(opts); i++ {
			opts[i] = r.impl.optSettings[i]
		}
		optCnt := container.NewBorder(r.impl.optSettingsName, nil, nil, nil, container.NewHScroll(container.NewHBox(opts...)))
		options = optCnt
	}

	// prepare profile panel
	profiles := make([]fyne.CanvasObject, len(r.impl.profiles))
	for i := 0; i < len(profiles); i++ {
		profiles[i] = r.impl.profiles[i]
	}

	// add prepared panels to content panel
	cnt := container.NewBorder(nil, options, nil, nil, container.NewHScroll(container.NewHBox(profiles...)))

	// set renderer to render new content panel
	r.WidgetRenderer = widget.NewSimpleRenderer(cnt)
}

type ParameterProfilePanel struct {
	widget.BaseWidget
	profiles        []*ParameterProfile
	optSettingsName *widget.Label
	optSettings     []*widget.Entry
	renderer        *ParameterProfilePanelRenderer
}

func (p *ParameterProfilePanel) CreateRenderer() fyne.WidgetRenderer {
	p.renderer = NewParameterProfilePanelRenderer(p)
	return p.renderer
}

func NewParameterProfilePanel(profiles ...*ParameterProfile) *ParameterProfilePanel {
	p := &ParameterProfilePanel{profiles: profiles, optSettingsName: widget.NewLabel("Settings"), optSettings: nil}
	p.ExtendBaseWidget(p)
	return p
}

func (p *ParameterProfilePanel) Add(profile *ParameterProfile) {
	p.profiles = append(p.profiles, profile)
	if p.renderer != nil {
		p.renderer.Update()
	}
}

func (p *ParameterProfilePanel) Remove(profile *ParameterProfile) {
	parameterIndex := slices.Index(p.profiles, profile)
	if parameterIndex == -1 {
		return
	}
	p.profiles = append(p.profiles[:parameterIndex], p.profiles[parameterIndex+1:]...)
	if p.renderer != nil {
		p.renderer.Update()
	}
}
