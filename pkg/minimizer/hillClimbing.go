package minimizer

import (
	"math"
	"slices"
)

type hillClimbingMinimizer struct {
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

func (h *hillClimbingMinimizer) Minimize(errFunc func([]float64) float64, x0 []float64, maxIterations int, minDelta float64) []float64 {
	bestEval := math.Inf(-1)
	currentNode := x0
	bestNode := x0
	for i := range maxIterations {
		valBuff := -errFunc(currentNode)
		if valBuff > bestEval {
			bestEval = valBuff
			bestNode = currentNode
		}
		L := permutations(currentNode, asymptoteGrowRate(float64(i)/float64(maxIterations), 10*minDelta, minDelta))
		tempMaxEval := math.Inf(-1)
		for _, x := range L {
			tmpValBuff := -errFunc(x)
			if tmpValBuff > tempMaxEval {
				currentNode = x
				tempMaxEval = tmpValBuff
			}
		}
		if slices.Equal(currentNode, bestNode) { // Stop if no neighbor is better => local maximum reached
			return bestNode
		}
	}
	return currentNode
}

func permutations(x0 []float64, dx float64) [][]float64 { //Simple version maybe update in future or use alternative minimizer
	res := make([][]float64, 2*len(x0))
	for i := 0; i < 2*len(x0); i++ {
		res[i] = make([]float64, len(x0))
		res[i] = slices.Clone(x0)
		res[i][i%len(x0)] += float64((i/len(x0))*2-1) * dx
	}
	return res
}

func asymptoteGrowRate(q, max, min float64) float64 {
	const qMapRate = 20
	return 1/(qMapRate*q+(1/(max-min))) + min
}
