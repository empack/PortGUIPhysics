package transformation

import (
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
	In  []T
	Wg  *sync.WaitGroup
	Out *K
}

func (b *BaseSegment[T, K]) Setup(in []T, wg *sync.WaitGroup, out *K) {
	b.In, b.Out, b.Wg = in, out, wg
}

func (b *BaseSegment[T, K]) Start(_ int) {
	b.Wg.Done()
}
