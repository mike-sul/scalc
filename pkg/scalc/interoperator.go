package scalc

import (
	"io"
)

const (
	InterOperatorId OperatorID = "INT"
)

type InterOperator struct {
	inputsets []SetReader
}

func NewInterOperator(inputsets []SetReader) SetReader {
	return NewProxy(&InterOperator{inputsets: inputsets})
}

func (i *InterOperator) ComputeNextValue() (int, error) {
	res := 0
	foundIntersection := false

	for {
		min, err := i.nextMinValue()
		if err != nil {
			return 0, err
		}

		foundIntersection = true

		for _, inputset := range i.inputsets {
			val, err := inputset.Peek()
			if err != nil {
				return 0, err
			}
			if val != min {
				// move next for the streams with minimum value
				for _, inputset1 := range i.inputsets {
					val, err := inputset1.Peek()
					if err == nil && val == min {
						inputset1.Next()
					}
				}
				foundIntersection = false
				break
			}
		}

		if foundIntersection {
			res = min
			break
		}
	}

	for _, inputset := range i.inputsets {
		val, err := inputset.Peek()
		if err == nil && val == res {
			// just move to the next item
			// EOF or any other error will be considered during next iteration
			inputset.Next()
		}
	}

	return res, nil
}

func (u *InterOperator) nextMinValue() (int, error) {
	min, err := u.getInitialMinValue()
	if err != nil {
		return 0, err
	}

	for _, inputset := range u.inputsets {
		val, err := inputset.Peek()
		if err != nil && err == io.EOF {
			continue
		}
		if err != nil {
			return 0, err
		}
		if val < min {
			min = val
		}
	}
	return min, nil
}

func (u *InterOperator) getInitialMinValue() (int, error) {
	for _, inputset := range u.inputsets {
		val, err := inputset.Peek()
		if err == nil {
			return val, err
		}
	}
	return 0, io.EOF
}
