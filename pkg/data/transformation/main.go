package transformation

import (
	"errors"
	"sync"
)

type PipelineSegment[T any, K any] interface {

	// Setup
	//
	// - in all outputs from the previous stage (read only)
	// - wg to keep track of all Parallel Segment states
	// - out pointer to the array element for output
	Setup(in []T, wg *sync.WaitGroup, out *K)

	// Start will be called parallel with the start function in all other Segment's in this Stage.
	//
	// - i will be the index of the segment in the stage
	//
	// **NOTE** the data from in should only be read and written only to out
	//
	// **IMPORTANT** After calculations are finished call wg.Done()
	Start(i int)
}

type BaseSegment[T any, K any] struct {
	in  []T
	wg  *sync.WaitGroup
	out *K
}

func (b *BaseSegment[T, K]) Setup(in []T, wg *sync.WaitGroup, out *K) {
	b.in, b.out, b.wg = in, out, wg
}

func (b *BaseSegment[T, K]) Start(_ int) {
	b.wg.Done()
	println(errors.New("forgot to override transformation.BaseSegment Start(i int) with own implement").Error())
}
