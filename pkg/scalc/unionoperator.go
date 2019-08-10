package scalc

import "io"

const (
	UnionOperatorId OperatorID = "SUM"
)

type UnionOperator struct {
	inputsets []SetReader
}

func NewUnionOperator(inputsets []SetReader) SetReader {
	return NewProxy(&UnionOperator{inputsets: inputsets})
}

func (u *UnionOperator) ComputeNextValue() (int, error) {
	min, error := u.getInitialMinValue()
	if error != nil {
		return 0, error
	}

	for _, inputset := range u.inputsets {

		val, err := inputset.Peek()
		// EOF of the ii/current Set has been reached, continue with the next one
		if err != nil && err == io.EOF {
			continue
		}
		// read/IO error has happened while peeking/reading a value from the set
		// so, just stop and return the error
		if err != nil {
			error = err
			break
		}

		if val < min {
			min = val
		}
	}

	for _, inputset := range u.inputsets {
		val, err := inputset.Peek()
		if err == nil && val == min {
			// just move to the next item
			// EOF or any other error will be considered during next iteration
			inputset.Next()
		}
	}

	return min, error
}

func (u *UnionOperator) getInitialMinValue() (int, error) {
	for _, inputset := range u.inputsets {
		val, err := inputset.Peek()
		if err == nil {
			return val, err
		}
	}
	return 0, io.EOF
}
