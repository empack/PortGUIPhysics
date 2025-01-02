package gui

import (
	"physicsGUI/pkg/data/transformation"
	"sync"
	"time"
)

type ScreenUpdater[T any, K any] struct {
	mu       sync.RWMutex
	dataIn   []T
	dataOut  []K
	pipeline *transformation.BasicAsyncPipeline[T, K]
	dirty    bool
	timer    chan time.Time
	ticker   *time.Ticker
	close    bool
}

func NewScreenUpdater[T any, K any](pipeline *transformation.BasicAsyncPipeline[T, K]) *ScreenUpdater[T, K] {
	u := &ScreenUpdater[T, K]{
		mu:       sync.RWMutex{},
		dataIn:   make([]T, 0),
		dataOut:  make([]K, 0),
		pipeline: pipeline,
		dirty:    true,
		timer:    make(chan time.Time, 2),
		ticker:   nil,
	}
	go u.updater()
	return u
}

func (u *ScreenUpdater[T, K]) updater() {
	isDirty := func() bool {
		u.mu.RLock()
		defer u.mu.RUnlock()
		return u.dirty
	}

	for {
		<-u.timer
		if !isDirty() {
			continue
		}
		u.mu.Lock()
		u.pipeline.Set(u.dataIn)
		u.dirty = false
		u.mu.Unlock()
		u.pipeline.Start()
		u.mu.Lock()
		u.dataOut = u.pipeline.Get()
		u.mu.Unlock()
	}
}

func (u *ScreenUpdater[T, K]) Loop(d time.Duration) {
	if u.ticker != nil {
		u.ticker.Reset(d)
	} else {
		u.ticker = time.NewTicker(d)
	}

	go func(c chan time.Time, ticker *time.Ticker) {
		for tick := range ticker.C {
			c <- tick
		}
	}(u.timer, u.ticker)
}

func (u *ScreenUpdater[T, K]) SetData(in ...T) {
	u.mu.Lock()
	u.dataIn = in
	u.dirty = true
	u.mu.Unlock()
}

func (u *ScreenUpdater[T, K]) GetData() []K {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.dataOut
}

func (u *ScreenUpdater[T, K]) Stop() {
	if u.ticker != nil {
		u.ticker.Stop()
		u.ticker = nil
	}
}

func (u *ScreenUpdater[T, K]) Destroy() {
	u.Stop()
	u.close = true
	close(u.timer)
}
