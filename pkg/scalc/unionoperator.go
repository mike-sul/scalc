package scalc

import (
	"io"
	"math"
)

const (
	UnionOperatorId OperatorID = "SUM"
)

type UnionOperator struct {
	inputSets []SetReader
}

func NewUnionOperator(sets []SetReader) SetReader {
	return NewProxy(&UnionOperator{inputSets: sets})
}

func (u *UnionOperator) ComputeNextValue() (int, error) {
	min, err := u.getNextMinValue()
	if err != nil {
		return min, err
	}

	// move to the next item in each set the top element of which equals to the computed next value (min)
	for _, inputSet := range u.inputSets {
		val, err := inputSet.Peek()
		if err == nil && val == min {
			// just move to the next item
			// EOF or any other error will be considered during next iteration
			inputSet.Next()
		}
	}

	return min, err
}

func (u *UnionOperator) getNextMinValue() (int, error) {
	var curMin = math.MaxInt64
	var retErr error = nil
	var nonEofSetCounter = 0

	for _, inputSet := range u.inputSets {
		val, err := inputSet.Peek()
		// read/IO error has happened while peeking/reading a value from the set
		// so, just stop and return the error
		if err != nil && err != io.EOF {
			retErr = err
			break
		}

		if err == io.EOF {
			continue
		}

		if val < curMin {
			curMin = val
		}
		nonEofSetCounter++
	}

	if nonEofSetCounter == 0 {
		retErr = io.EOF
	}
	return curMin, retErr
}