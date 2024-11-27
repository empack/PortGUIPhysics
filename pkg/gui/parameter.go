package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"maps"
	"physicsGUI/pkg/util/option"
	"slices"
	"strconv"
	"strings"
)

const (
	ProfileDefaultEdensityID  = iota
	ProfileDefaultRoughnessID = iota
	ProfileDefaultThicknessID = iota
)

type FilteredEntry struct {
	widget.Entry
	allowedTyping []rune
}

// NewFilteredEntry creates a widget.Entry, that only accept the keys specified in allowedRunes
func NewFilteredEntry(allowedRunes []rune) *FilteredEntry {
	entry := &FilteredEntry{
		allowedTyping: allowedRunes,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *FilteredEntry) TypedRune(r rune) {
	if slices.Contains(e.allowedTyping, r) {
		e.Entry.TypedRune(r)
	}
}

type Parameter struct {
	widget.BaseWidget
	name         *widget.Label
	defaultValue float64
	valEntry     *FilteredEntry
	minEntry     *FilteredEntry
	maxEntry     *FilteredEntry
	locked       *widget.Check
	OnChanged    func()
	maxSize      fyne.Size
	minSize      fyne.Size
	objects      []fyne.CanvasObject // weil im renderer unterschiedliche Structs für rendern und Layout verwendet werden warum auch immer
}

func NewParameter(name string, defaultValue float64) *Parameter {
	numberRunes := []rune("0123456789.e-+")

	p := Parameter{
		name:         widget.NewLabel(name),
		defaultValue: defaultValue,
		valEntry:     NewFilteredEntry(numberRunes),
		minEntry:     NewFilteredEntry(numberRunes),
		maxEntry:     NewFilteredEntry(numberRunes),
		locked:       widget.NewCheck("", nil),
		OnChanged:    func() {},
		minSize:      fyne.NewSize(240, 80),
		maxSize:      fyne.NewSize(400, 200),
	}
	p.locked.OnChanged = func(_ bool) {
		p.OnChanged()
	}
	p.valEntry.OnChanged = func(_ string) {
		p.OnChanged()
	}
	p.minEntry.OnChanged = func(_ string) {
		p.OnChanged()
	}
	p.maxEntry.OnChanged = func(_ string) {
		p.OnChanged()
	}
	return &p
}
func (p *Parameter) Refresh() {
	p.valEntry.Refresh()
	p.minEntry.Refresh()
	p.maxEntry.Refresh()
	p.locked.Refresh()
	p.BaseWidget.Refresh()
}

func (this *Parameter) MinSize() fyne.Size {
	altMinX := max(this.name.MinSize().Width,
		this.valEntry.MinSize().Width+
			this.locked.MinSize().Width+
			max(this.minEntry.MinSize().Width, this.maxEntry.MinSize().Width))
	altMinY := this.name.MinSize().Height +
		max(this.valEntry.MinSize().Height,
			this.locked.MinSize().Height,
			max(this.minEntry.MinSize().Height+this.maxEntry.MinSize().Height))
	return fyne.NewSize(max(altMinX, this.minSize.Width), max(altMinY, this.minSize.Height))
}

// Clear removes all content from user input fields and Refreshes
//
// **note** this does not affect the widget.Check for IsFixed()
func (this *Parameter) Clear() {
	if this.valEntry != nil {
		this.valEntry.Text = ""
	}
	if this.minEntry != nil {
		this.minEntry.Text = ""
	}
	if this.maxEntry != nil {
		this.maxEntry.Text = ""
	}
	this.Refresh()
}

// GetValue return the value in the value field
//
// - when input is empty or could not be parsed this.defaultValue is returned instead
// - when input could not be parsed error contains the error to display to user
// - else returns parsed value of valEntry field
func (this *Parameter) GetValue() (float64, error) {
	if this.valEntry.Text == "" {
		this.valEntry.SetText(fmt.Sprintf("%f", this.defaultValue))
		return this.defaultValue, nil
	}
	val, err := strconv.ParseFloat(this.valEntry.Text, 64)
	if err != nil {
		return this.defaultValue, errors.New("Float_Parsing_Error: Error while parsing Value input to float please adjust input of all marked fields")
	}
	return val, nil
}

// GetMin return the value in the Min field
//
// - when input is empty or could not be parsed option.None is returned instead
// - when input could not be parsed error contains the error to display to user
// - else returns parsed value of minEntry field wrapped in option.Some
func (this *Parameter) GetMin() (option.Option[float64], error) {
	if this.valEntry.Text == "" {
		return option.None[float64](), nil
	}
	val, err := strconv.ParseFloat(this.valEntry.Text, 64)
	if err != nil {
		return option.None[float64](), errors.New("Float_Parsing_Error: Error while parsing Min input to float please adjust input of all marked fields")
	}
	return option.Some[float64](&val), nil // Maybe change for better memory layout
}

// GetMax return the value in the Max field
//
// - when input is empty or could not be parsed option.None is returned instead
// - when input could not be parsed error contains the error to display to user
// - else returns parsed value of maxEntry field wrapped in option.Some
func (this *Parameter) GetMax() (option.Option[float64], error) {
	if this.valEntry == nil || this.valEntry.Text == "" {
		return option.None[float64](), nil
	}
	val, err := strconv.ParseFloat(this.valEntry.Text, 64)
	if err != nil {
		return option.None[float64](), errors.New("Float_Parsing_Error: Error while parsing Max input to float please adjust input of all marked fields")
	}
	return option.Some[float64](&val), nil // Maybe change for better memory layout
}
func (this *Parameter) IsFixed() option.Option[bool] {
	if this.locked != nil {
		return option.Some[bool](&this.locked.Checked)
	}
	return option.None[bool]()
}

func (this *Parameter) CreateRenderer() fyne.WidgetRenderer {
	parsable := func(s string) error {
		if s == "" {
			return nil
		}
		_, err := strconv.ParseFloat(s, 64)
		return err
	}
	this.valEntry.MultiLine = false
	this.valEntry.Validator = parsable
	this.valEntry.PlaceHolder = fmt.Sprintf("%f", this.defaultValue)
	this.valEntry.Scroll = container.ScrollNone
	this.valEntry.Refresh()
	if this.maxEntry != nil {
		this.maxEntry.MultiLine = false
		this.maxEntry.Validator = parsable
		this.maxEntry.PlaceHolder = "Max"
		this.maxEntry.Scroll = container.ScrollNone
		this.maxEntry.Refresh()
	}
	if this.minEntry != nil {
		this.minEntry.MultiLine = false
		this.minEntry.Validator = parsable
		this.minEntry.PlaceHolder = "Min"
		this.minEntry.Scroll = container.ScrollNone
		this.minEntry.Refresh()
	}
	if this.maxEntry != nil && this.minEntry != nil {
		return widget.NewSimpleRenderer(container.NewVBox(this.name, container.NewHBox(this.locked, this.valEntry, container.NewVBox(this.maxEntry, this.minEntry))))
	}
	return widget.NewSimpleRenderer(container.NewVBox(this.name, container.NewHBox(this.locked, this.valEntry)))
}

type Profile struct {
	widget.BaseWidget
	name      *widget.Entry
	removeBtn *widget.Button
	idStart   int
	parameter map[int]*Parameter
}

// NewBlankProfile creates a new Profile with nothing but a name
func NewBlankProfile(name string) *Profile {
	p := &Profile{
		name:      widget.NewEntry(),
		removeBtn: nil,
		idStart:   0,
		parameter: map[int]*Parameter{},
	}
	p.name.SetText(name)
	p.ExtendBaseWidget(p)
	return p
}
func (this *Profile) SetParameter(id int, parameter *Parameter) {
	if id > this.idStart {
		this.idStart = id
	}
	this.parameter[id] = parameter
}
func (this *Profile) AddParameter(parameter *Parameter) int {
	newID := this.idStart + 1
	this.parameter[newID] = parameter
	this.idStart = newID
	this.Refresh()
	return newID
}
func (this *Profile) RemoveParameter(id int) {
	this.parameter[id] = nil
	keys := maps.Keys(this.parameter)
	this.idStart = slices.Max(slices.Collect(keys))

	this.Refresh()
}
func NewDefaultProfile(name string, roughnessName string, defaultRoughness float64, edensityName string, defaultEdensity float64, thicknessName string, defaultThickness float64) *Profile {
	parameter := make(map[int]*Parameter, 3)
	parameter[ProfileDefaultRoughnessID] = NewParameter(roughnessName, defaultRoughness)
	parameter[ProfileDefaultEdensityID] = NewParameter(edensityName, defaultEdensity)
	parameter[ProfileDefaultThicknessID] = NewParameter(thicknessName, defaultThickness)
	profile := &Profile{
		name:      widget.NewEntry(),
		removeBtn: widget.NewButton("🗑", func() {}),
		idStart:   max(ProfileDefaultThicknessID, ProfileDefaultEdensityID, ProfileDefaultRoughnessID),
		parameter: parameter,
	}
	// Default button function clears inputs
	profile.removeBtn = widget.NewButton("🗑", func() {
		profile.Clear()
	})
	profile.name.Text = name
	profile.ExtendBaseWidget(profile)
	return profile
}
func (this *Profile) Refresh() {
	this.BaseWidget.Refresh()
}

func (this *Profile) CreateRenderer() fyne.WidgetRenderer {
	var obj []fyne.CanvasObject
	for v := range maps.Values(this.parameter) {
		if v != nil {
			obj = append(obj, v)
		}
	}
	var cnt = container.NewBorder(nil, nil, nil, nil, this.name)
	if this.removeBtn != nil {
		cnt = container.NewBorder(nil, nil, nil, this.removeBtn, this.name)
	}
	return widget.NewSimpleRenderer(container.NewVBox(append([]fyne.CanvasObject{cnt}, obj...)...))
}

// Clear removes calls Parameter.Clear() on all Parameter's and Refreshes afterward
func (this *Profile) Clear() {
	for _, parameter := range this.parameter {
		if parameter != nil {
			parameter.Clear()
		}
	}
	this.Refresh()
}

type ProfilePanel struct {
	widget.BaseWidget
	base      *Profile
	bulk      *Profile
	profiles  []*Profile
	addButton *widget.Button
	renderer  *ProfilePanelRenderer
}

func (this *ProfilePanel) Resize(size fyne.Size) {
	if this.renderer != nil {
		this.renderer.Layout(size)
	}
	this.BaseWidget.Resize(size)
}

// NewProfilePanel creates a new ProfilePanel with given Profile's
//
// # If the given Profiles are empty one default Profile element gets added
//
// The ProfilePanel includes a Base Profile at the start and a Bulk Profile (without ProfileDefaultRoughnessID Parameter) at the end, as well as a add button to create new Profile's
func NewProfilePanel(profiles ...*Profile) *ProfilePanel {
	//TODO read default values from some settings file
	defaultRoughness := 10.0
	defaultEdensity := 1.0
	defaultThickness := 1.0

	base := NewDefaultProfile("Base", "Roughness Base/Slab1", defaultRoughness, "Edensity Base", defaultEdensity, "Thickness Base", defaultThickness)
	bulk := NewDefaultProfile("Bulk", "Roughness Bulk", 0.0, "Edensity Bulk", defaultEdensity, "Thickness Bulk", defaultThickness)
	bulk.parameter[ProfileDefaultRoughnessID] = nil

	p := &ProfilePanel{
		base:     base,
		bulk:     bulk,
		profiles: profiles,
	}
	if profiles == nil || len(profiles) == 0 {
		profileName := fmt.Sprintf("Slab 1")
		roughnessName := fmt.Sprintf("Rougthness Slab 1/Bulk")
		edensityName := fmt.Sprintf("Edensity Slab 1")
		thicknessName := fmt.Sprintf("Thickness Slab 1")
		baseProfile := NewDefaultProfile(profileName, roughnessName, defaultRoughness, edensityName, defaultEdensity, thicknessName, defaultThickness)
		p.AddProfile(baseProfile)
	}

	p.addButton = widget.NewButton("+", func() {
		profileName := fmt.Sprintf("Slab %d", len(p.profiles)+1)
		roughnessName := fmt.Sprintf("Rougthness Slab %d/Bulk", len(p.profiles)+1)
		edensityName := fmt.Sprintf("Edensity Slab %d", len(p.profiles)+1)
		thicknessName := fmt.Sprintf("Thickness Slab %d", len(p.profiles)+1)
		newP := NewDefaultProfile(profileName, roughnessName, defaultRoughness, edensityName, defaultEdensity, thicknessName, defaultThickness)
		p.AddProfile(newP) //TODO add Settings for default values
	})
	return p
}

// AddProfile adds the given Profile to the Panel.
//
// - The name of the ProfileDefaultRoughnessID Parameter will be updated, if the Parameter exists in the Profile,
// as well as the name of the ProfileDefaultRoughnessID Parameter of the previous Element
//
// - Adds a remove button to the Profile
func (this *ProfilePanel) AddProfile(profile *Profile) {
	if len(this.profiles) > 0 {
		param := this.profiles[len(this.profiles)-1].parameter[ProfileDefaultRoughnessID]
		if param != nil {
			param.name.SetText(fmt.Sprintf("Rougthness Slab %d/Slab %d", len(this.profiles), len(this.profiles)+1))
		}
	}
	this.profiles = append(this.profiles, profile)
	profile.removeBtn = widget.NewButton("🗑", func() {
		this.RemoveProfile(profile)
	})
	if this.renderer != nil {
		this.renderer.Update()
	}
}

// RemoveProfile removes the given profile address from the panel and updates the names of the parameters from the other profiles
//
// - The numbers in layers with custom names are also changed, when they match with the layer number
// Example: Layer with name 'Layer 2' is the second layer and becomes first layer, it's name gets updates to 'Layer 1'
//
// **Note** the last Profile can not be removed, it calls Profile.Clear() instead and resets the name to a generic "Slab 1"
func (this *ProfilePanel) RemoveProfile(profile *Profile) {
	i := slices.Index(this.profiles, profile)
	if i >= 0 {
		if len(this.profiles) > 1 {
			if i != len(this.profiles)-1 {
				for j := i + 1; j < len(this.profiles); j++ {
					roughnessName := fmt.Sprintf("Rougthness Slab %d/%d", j, j+1)
					edensityName := fmt.Sprintf("Edensity Slab %d", j)
					thicknessName := fmt.Sprintf("Thickness Slab %d", j)
					numberIndex := strings.LastIndex(this.profiles[j].name.Text, fmt.Sprint(j+1))
					if numberIndex != -1 {
						name := this.profiles[j].name.Text
						newNumber := fmt.Sprint(j)
						this.profiles[j].name.SetText(name[:numberIndex] + newNumber + name[numberIndex+len(newNumber):])
					}
					if this.profiles[j].parameter[ProfileDefaultRoughnessID] != nil {
						this.profiles[j].parameter[ProfileDefaultRoughnessID].name.SetText(roughnessName)
					}
					if this.profiles[j].parameter[ProfileDefaultEdensityID] != nil {
						this.profiles[j].parameter[ProfileDefaultEdensityID].name.SetText(edensityName)
					}
					if this.profiles[j].parameter[ProfileDefaultThicknessID] != nil {
						this.profiles[j].parameter[ProfileDefaultThicknessID].name.SetText(thicknessName)
					}
				}
				this.profiles = append(this.profiles[:i], this.profiles[i+1:]...)
			} else {
				this.profiles = this.profiles[:i]
			}
			if this.profiles[len(this.profiles)-1].parameter[ProfileDefaultRoughnessID] != nil {
				this.profiles[len(this.profiles)-1].parameter[ProfileDefaultRoughnessID].name.SetText(fmt.Sprintf("Rougthness Slab %d/Bulk", len(this.profiles)))
			}
		} else {
			this.profiles[i].Clear()
			this.profiles[i].name.SetText("Slab 1")
			this.profiles[i].Refresh()
		}
	}
	this.renderer.Update()
}

func (this *ProfilePanel) CreateRenderer() fyne.WidgetRenderer {
	this.renderer = NewProfilePanelRenderer(this)
	return this.renderer
}

type ProfilePanelRenderer struct {
	obj    *ProfilePanel
	layout fyne.WidgetRenderer
}

// Update forces Recalculation of all Profiles on this panel
//
// **Note** Call this, when visual components where added or removed
func (p *ProfilePanelRenderer) Update() {
	objects := make([]fyne.CanvasObject, len(p.obj.profiles)+3)
	for i, profile := range p.obj.profiles {
		objects[i+1] = profile
	}
	objects[0] = p.obj.base
	objects[len(objects)-2] = p.obj.addButton
	objects[len(objects)-1] = p.obj.bulk

	center := container.NewHBox(objects...)
	cnt := container.NewHScroll(center)
	p.layout = widget.NewSimpleRenderer(cnt)
	p.Refresh()
}

func (p *ProfilePanelRenderer) Destroy() {}

func (p *ProfilePanelRenderer) Layout(size fyne.Size) {
	p.layout.Layout(size)
}

func (p *ProfilePanelRenderer) MinSize() fyne.Size {
	return p.layout.MinSize()
}

func (p *ProfilePanelRenderer) Objects() []fyne.CanvasObject {
	if p.layout == nil {
		p.Update()
	}
	return p.layout.Objects()
}

func (p *ProfilePanelRenderer) Refresh() {
	p.layout.Layout(p.obj.Size())
	p.layout.Refresh()
}

func NewProfilePanelRenderer(obj *ProfilePanel) *ProfilePanelRenderer {
	return &ProfilePanelRenderer{
		obj:    obj,
		layout: nil,
	}
}
