package data

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

type ParameterHandler struct {
	ChangeListenerGroup[*Parameter]
	parameters   map[ParameterID]*Parameter
	uidGenerator *ParameterUIDGenerator
}

func NewParameterHandler() *ParameterHandler {
	uidGenerator := NewParameterUIDGenerator()
	return &ParameterHandler{
		parameters:   make(map[ParameterID]*Parameter),
		uidGenerator: uidGenerator,
	}
}

func (p *ParameterHandler) Add(parameter *Parameter) ParameterID {
	request := parameter.uid
	parameter.uid = "unregistered"
	p.ChangeID(parameter, request)
	p.parameters[parameter.uid] = parameter
	p.trigger(nil, parameter)
	return parameter.uid
}

func (p *ParameterHandler) Remove(parameter *Parameter) {
	if param, ok := p.parameters[parameter.uid]; !ok {
		return
	} else {
		p.trigger(param, nil)
		delete(p.parameters, parameter.uid)
		p.uidGenerator.unregister(parameter.uid)
	}
}

func (p *ParameterHandler) TryChangeID(parameter *Parameter, requestedID ParameterID) bool {
	b := p.uidGenerator.register(ParameterID(strings.ToLower(string(requestedID))))
	if !b {
		return false
	} else {
		oldID := parameter.uid
		parameter.uid = requestedID
		if oldID != parameter.uid {
			p.uidGenerator.unregister(oldID)
		}
		return true
	}
}

func (p *ParameterHandler) ChangeID(parameter *Parameter, requestedID ParameterID) {
	b := p.TryChangeID(parameter, requestedID)
	if !b {
		if parameter.uid != requestedID {
			p.uidGenerator.unregister(parameter.uid)
		}
		idBase := "Parameter"
		if parameter.class != "" {
			idBase = string(parameter.class)
		}
		newID := p.uidGenerator.generateID(idBase)
		parameter.uid = newID
	}
}

func (p *ParameterHandler) GetByUID(uid ParameterID) *Parameter {
	if parameter, ok := p.parameters[ParameterID(strings.ToLower(string(uid)))]; ok {
		return parameter
	} else {
		return nil
	}
}

func (p *ParameterHandler) GetByClass(classname string) []*Parameter {
	res := make([]*Parameter, 0)
	for _, param := range p.parameters {
		if strings.ToLower(string(param.class)) == strings.ToLower(classname) {
			res = append(res, param)
		}
	}
	slices.SortFunc(res, func(a, b *Parameter) int {
		maxID := fmt.Sprint("-", math.MaxInt)
		aCleared := strings.Replace(strings.Replace(string(a.uid), "-a", "-0", 1), "-b", maxID, 1)
		bCleared := strings.Replace(strings.Replace(string(b.uid), "-a", "-0", 1), "-b", maxID, 1)
		return strings.Compare(aCleared, bCleared)
	})
	return res
}

func (p *ParameterHandler) GetAll() []*Parameter {
	// maps.Values(p.parameters)
	values := make([]*Parameter, len(p.parameters))
	index := 0
	for _, v := range p.parameters {
		values[index] = v
		index++
	}
	slices.SortFunc(values, func(a, b *Parameter) int {
		return strings.Compare(string(a.uid), string(b.uid))
	})
	return values
}

type ParameterUIDGenerator struct {
	ids map[ParameterID]bool
}

func NewParameterUIDGenerator() *ParameterUIDGenerator {
	return &ParameterUIDGenerator{
		ids: make(map[ParameterID]bool),
	}
}

func (g *ParameterUIDGenerator) unregister(uid ParameterID) {
	delete(g.ids, uid)
}

func (g *ParameterUIDGenerator) register(uid ParameterID) bool {
	if _, ok := g.ids[uid]; ok {
		return false
	} else {
		g.ids[uid] = true
		return true
	}
}

func (g *ParameterUIDGenerator) generateID(base string) ParameterID {
	if base == "" {
		base = "parameter"
	}
	guessBase := strings.ToLower(base) + "-"
	guess := ParameterID(guessBase)
	for ok, i := false, 1; i <= math.MaxInt && !ok; ok, i = g.register(guess), i+1 {
		guess = ParameterID(guessBase + fmt.Sprint(i))
	}
	// assert we have more integers than current parameters
	return guess
}
