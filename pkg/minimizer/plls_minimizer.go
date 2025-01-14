package minimizer

import (
	"slices"
	"sync"
)

type parallelLinearLocalSearch[T Number] struct {
	minDelta T
}

func (p *parallelLinearLocalSearch[T]) Minimize(problem *AsyncMinimiserProblem[T]) {
	// Get information to set up minimisation
	wg := new(sync.WaitGroup)
	problem.lock.RLock()
	idCount := len(problem.parameter)
	problem.lock.RUnlock()

	// setup minimisation
	wg.Add(idCount)

	// spawn parallel worker routines
	for i := 0; i < idCount; i++ {
		go p.plsWorker(i, wg, problem)
	}

	// wait for completion
	wg.Wait()

	// mark minimization as done
	problem.lock.Lock()
	problem.config.loopCount = 0
	problem.lock.Unlock()
}

func (p *parallelLinearLocalSearch[T]) plsWorker(id int, wg *sync.WaitGroup, problem *AsyncMinimiserProblem[T]) {
	// read config and make local of necessary data copy
	problem.lock.RLock()
	useLock := problem.config.parallelReads
	minDelta := p.minDelta
	maxIterations := problem.config.loopCount
	parameters := problem.parameter
	minv := problem.minima[id]
	maxv := problem.maxima[id]
	problem.lock.RUnlock()

	// Setup
	innerWg := new(sync.WaitGroup)

	// run minimisation
	for i := 0; i < maxIterations; i++ {
		if useLock {
			problem.lock.RLock()
			parameters = problem.parameter
			problem.lock.RUnlock()
		}
		var parameter T
		var errP T
		var errM T
		guessP := slices.Clone(parameters)
		guessP[id] = min(guessP[id]+minDelta, maxv)
		guessM := slices.Clone(parameters)
		guessM[id] = max(guessM[id]-minDelta, minv)

		innerWg.Add(2)

		go func() {
			errP = problem.errorFunction(guessP)
			innerWg.Done()
		}()
		go func() {
			errM = problem.errorFunction(guessM)
			innerWg.Done()
		}()

		// wait for error function to finish
		innerWg.Wait()

		// go in direction where error goes smaller
		if errP < errM {
			parameter = guessP[id]
		} else {
			parameter = guessM[id]
		}

		// check if minima was found / no better option was found
		if parameter == parameters[id] {
			break
		}

		// write back current result
		if useLock {
			problem.lock.Lock()
			problem.parameter[id] = parameter
			problem.lock.Unlock()
		} else {
			problem.parameter[id] = parameter
		}
	}

	// mark parameter as minimized
	wg.Done()
}
