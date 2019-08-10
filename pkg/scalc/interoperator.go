package scalc

import (
	"math"
)

const (
	InterOperatorId OperatorID = "INT"
)

type InterOperator struct {
	inputSets []SetReader
}

func NewInterOperator(sets []SetReader) SetReader {
	//return NewProxy(&InterOperator{inputSets: sets})
	return NewChannelProxy(&InterOperator{inputSets: sets})
}

func (i *InterOperator) ComputeNextValue() (int, error) {
	var retErr error = nil
	var min = math.MaxInt64

	for retErr == nil {

		min = math.MaxInt64
		// find a min value across the top elements of the input sets
		for _, inputSet := range i.inputSets {
			val, err := inputSet.Peek()
			// read/IO error has happened while peeking/reading a value from the set
			// so, just stop and return the error
			// also, if io.EOF has been reached for at least one set then stop
			if err != nil {
				retErr = err
				break
			}

			if val < min {
				min = val
			}

		}

		// move to the next item in the sets where the top equals to the min
		// if each set top element equals to the min then an intersection is found
		foundIntersection := true
		for _, inputSet := range i.inputSets {
			val, err := inputSet.Peek()
			if err != nil {
				retErr = err
				break
			}
			if val == min {
				// just move to the next item
				// EOF or any other error will be considered during next iteration
				inputSet.Next()
			} else {
				foundIntersection = false
			}
		}

		if foundIntersection { break }
	}

	return min, retErr
}
