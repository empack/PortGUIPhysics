package minimizer

import (
	"errors"
	"sync"
)

type Number interface {
	~uint8 | ~uint32 | ~uint64 |
		~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

type MinimiserConfig struct {
	loopCount     int
	parallelReads bool
}

func NewParallelConfig(maxLoopCount int) *MinimiserConfig {
	return &MinimiserConfig{
		loopCount:     maxLoopCount,
		parallelReads: true,
	}
}
func NewAsyncConfig(maxLoopCount int) *MinimiserConfig {
	return &MinimiserConfig{
		loopCount:     maxLoopCount,
		parallelReads: false,
	}
}

type AsyncMinimiserProblem[T Number] struct {
	lock          sync.RWMutex
	config        *MinimiserConfig
	parameter     []T
	minima        []T
	maxima        []T
	errorFunction func(parameter []T) T
}

func NewProblem[T Number](x0, minima, maxima []T, errorFunction func(parameter []T) T, config *MinimiserConfig) *AsyncMinimiserProblem[T] {
	return &AsyncMinimiserProblem[T]{
		lock:          sync.RWMutex{},
		config:        config,
		parameter:     x0,
		minima:        minima,
		maxima:        maxima,
		errorFunction: errorFunction,
	}
}

func (p *AsyncMinimiserProblem[T]) GetCurrentParameters() ([]T, error) {
	if p.config.parallelReads {
		p.lock.RLock()
		defer p.lock.RUnlock()
		return p.parameter, nil
	} else if p.config.loopCount == 0 {
		return p.parameter, nil
	} else {
		return nil, errors.New("can not get parameters while minimizing with without parallel read option")
	}
}

// Pause try to pause a running problem
//
// WARNING:
// Ensure, that this problem was not paused previously without successful resume otherwise it can crash you program with fatal error
func (p *AsyncMinimiserProblem[_]) Pause() error {
	if !p.config.parallelReads {
		return errors.New("can not pause with without parallel read option")
	} else if p.config.loopCount == 0 {
		return errors.New("can not pause what is already completed")
	}
	p.lock.Lock()
	return nil
}

// Resume try to resume a paused problem
//
// WARNING:
// Ensure, that you have successfully paused this problem previously otherwise it can crash you program with fatal error
func (p *AsyncMinimiserProblem[_]) Resume() error {
	if !p.config.parallelReads {
		return errors.New("can not resume without parallel read option")
	} else if p.config.loopCount == 0 {
		return errors.New("can not resume what is already completed")
	}
	p.lock.Unlock()
	return nil
}

type Minimizer[T Number] interface {
	Minimize(problem *AsyncMinimiserProblem[T])
}

var (
	FloatMinimizerPLLS Minimizer[float64] = &parallelLinearLocalSearch[float64]{minDelta: 1e-6}
	IntMinimizerPLLS   Minimizer[int64]   = &parallelLinearLocalSearch[int64]{minDelta: 1}
	FloatMinimizerHC   Minimizer[float64] = &hillClimbingMinimizer[float64]{minDelta: 1e-6}
	IntMinimizerHC     Minimizer[int64]   = &hillClimbingMinimizer[int64]{minDelta: 1}
)
