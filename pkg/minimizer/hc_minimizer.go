package minimizer

import (
	"slices"
	"sync"
)

type hillClimbingMinimizer[T Number] struct {
	minDelta T
}

/*
Algo (Hill Climbing):
    bestEval = -INF
    currentNode = startNode
    bestNode = None
    for MAX times:
        if EVAL(currentNode) > bestEval:
            bestEval = EVAL(currentNode)
            bestNode = currentNode
        L = NEIGHBORS(currentNode)
        tempMaxEval = -INF
        for all x in L:
            if EVAL(x) > tempMaxEval:
                currentNode = x
                tempMaxEval = EVAL(x)
    return currentNode
*/

func (h *hillClimbingMinimizer[T]) Minimize(problem *AsyncMinimiserProblem[T]) {
	problem.lock.RLock()
	useLocks := problem.config.parallelReads
	bestNode := problem.parameter
	minv := problem.minima
	maxv := problem.maxima
	problem.lock.RUnlock()

	parameterCount := len(bestNode)
	wg := new(sync.WaitGroup)

	shouldLoop := func() bool {
		if useLocks {
			problem.lock.RLock()
			defer problem.lock.RUnlock()
			return problem.config.loopCount != 0
		} else {
			return problem.config.loopCount != 0
		}
	}
	bestEval := problem.errorFunction(bestNode)
	for shouldLoop() {
		// calculate errors of neighbors
		neighborErrors := make([]T, 2*parameterCount)
		wg.Add(2 * parameterCount)
		for i := 0; i < parameterCount; i++ {
			go func(id int) {
				dir := slices.Clone(bestNode)
				dir[id] = min(dir[id]+h.minDelta, maxv[id])
				neighborErrors[id] = problem.errorFunction(dir)
				wg.Done()
			}(i)
			go func(id int) {
				dir := slices.Clone(bestNode)
				dir[id] = max(dir[id]-h.minDelta, minv[id])
				neighborErrors[parameterCount+id] = problem.errorFunction(dir)
				wg.Done()
			}(i)
		}
		wg.Wait()

		// find neighbor node with the lowest error
		mini := 0
		for i := 0; i < len(neighborErrors); i++ {
			if neighborErrors[i] < neighborErrors[mini] {
				mini = i
			}
		}

		// check if local minima was reached
		if neighborErrors[mini] >= bestEval {
			problem.lock.Lock()
			problem.config.loopCount = 0
			problem.lock.Unlock()
			continue
		}

		if mini < parameterCount {
			bestNode[mini] += h.minDelta
		} else {
			bestNode[mini-parameterCount] -= h.minDelta
		}
		bestEval = neighborErrors[mini]

		problem.lock.Lock()
		problem.config.loopCount--
		problem.parameter = bestNode
		problem.lock.Unlock()
	}
}
