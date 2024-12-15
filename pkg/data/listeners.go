package data

import (
	"fyne.io/fyne/v2/data/binding"
	"slices"
)

type ListenerGroup struct {
	listener []binding.DataListener
}

func (p *ListenerGroup) DataChanged() {
	p.trigger()
}
func (p *ListenerGroup) AddListener(listener binding.DataListener) {
	p.listener = append(p.listener, listener)
}
func (p *ListenerGroup) RemoveListener(listener binding.DataListener) {
	index := slices.Index(p.listener, listener)
	if index != -1 {
		p.listener = append(p.listener[:index], p.listener[index+1:]...)
	}
}
func (p *ListenerGroup) trigger() {
	for _, listener := range p.listener {
		listener.DataChanged()
	}
}
func NewListenerGroup() *ListenerGroup {
	return &ListenerGroup{}
}

type ChangeListener[T any] interface {
	DataChanged(T, T)
}
type SimpleChangeListener[T any] struct {
	callback func(T, T)
}

func (s *SimpleChangeListener[T]) DataChanged(old T, new T) {
	s.callback(old, new)
}
func NewChangeListener[T any](callback func(T, T)) *SimpleChangeListener[T] {
	return &SimpleChangeListener[T]{callback}
}

type ChangeListenerGroup[T any] struct {
	listener []ChangeListener[T]
}

func (p *ChangeListenerGroup[T]) DataChanged(old T, new T) {
	p.trigger(old, new)
}
func (p *ChangeListenerGroup[T]) AddListener(listener ChangeListener[T]) {
	p.listener = append(p.listener, listener)
}
func (p *ChangeListenerGroup[T]) RemoveListener(listener ChangeListener[T]) {
	index := slices.Index(p.listener, listener)
	if index != -1 {
		p.listener = append(p.listener[:index], p.listener[index+1:]...)
	}
}
func (p *ChangeListenerGroup[T]) HasListener(listener ChangeListener[T]) bool {
	return slices.Index(p.listener, listener) != -1
}
func (p *ChangeListenerGroup[T]) trigger(old T, new T) {
	for _, listener := range p.listener {
		listener.DataChanged(old, new)
	}
}
func NewChangeListenerGroup[T any]() *ChangeListenerGroup[T] {
	return &ChangeListenerGroup[T]{}
}
