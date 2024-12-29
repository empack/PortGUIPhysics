package transformation

import (
	"slices"
	"sync"
)

type PipelineStage interface {
	Set(in []any)
	Start()
	Get() []any
}

type Stage[T any, K any] struct {
	wg       sync.WaitGroup
	segments []PipelineSegment[T, K]
	input    []T
	output   []K
}

func NewStage[T any, K any](segments ...PipelineSegment[T, K]) *Stage[T, K] {
	return &Stage[T, K]{
		wg:       sync.WaitGroup{},
		segments: segments,
		input:    nil,
		output:   nil,
	}
}
func (s *Stage[T, K]) Set(in []any) {
	s.input = make([]T, len(in))
	for i := 0; i < len(in); i++ {
		s.input[i] = in[i].(T)
	}
}
func (s *Stage[T, K]) Start() {
	s.wg.Add(len(s.segments))
	s.output = make([]K, len(s.segments))
	for i, segment := range s.segments {
		segment.Setup(s.input, &s.wg, &s.output[i])
		go segment.Start(i)
	}
	s.wg.Wait()
}
func (s *Stage[T, K]) Get() []any {
	output := make([]any, len(s.segments))
	for i := 0; i < len(s.output); i++ {
		output[i] = s.output[i]
	}
	return output
}

func (s *Stage[T, K]) AddSegment(segment PipelineSegment[T, K]) {
	s.segments = append(s.segments, segment)
}
func (s *Stage[T, K]) RemoveSegment(segment PipelineSegment[T, K]) {
	if i := slices.Index(s.segments, segment); i != -1 {
		s.segments = append(s.segments[:i], s.segments[i+1:]...)
	}
}

type BasicAsyncPipeline[T any, K any] struct {
	input  []T
	output []K
	stages []PipelineStage
}

func NewBasicAsyncPipeline[T any, K any]() *BasicAsyncPipeline[T, K] {
	return &BasicAsyncPipeline[T, K]{}
}

func (p *BasicAsyncPipeline[T, K]) Set(in []T) {
	p.input = in
}
func (p *BasicAsyncPipeline[T, K]) Start() {
	in := make([]any, len(p.input))
	for i := 0; i < len(p.input); i++ {
		in[i] = p.input[i]
	}

	for i := 0; i < len(p.stages); i++ {
		p.stages[i].Set(in)
		p.stages[i].Start()
		in = p.stages[i].Get()
	}
	p.output = make([]K, len(in))
	for i := 0; i < len(p.output); i++ {
		p.output[i] = in[i].(K)
	}
}
func (p *BasicAsyncPipeline[T, K]) Get() []K {
	return p.output
}

func (p *BasicAsyncPipeline[T, K]) AddStage(stage PipelineStage) {
	p.stages = append(p.stages, stage)
}
func (p *BasicAsyncPipeline[T, K]) RemoveStage(stage PipelineStage) {
	if i := slices.Index(p.stages, stage); i != -1 {
		p.stages = append(p.stages[:i], p.stages[i+1:]...)
	}
}
