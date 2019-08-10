package scalc

import (
	"io"
	"math"
)

const (
	DifOperatorId OperatorID = "DIF"
)

type DiffOperator struct {
	inputSets []SetReader
}

func NewDiffOperator(sets []SetReader) SetReader {
	return NewProxy(&DiffOperator{inputSets: sets})
}

func (d *DiffOperator) ComputeNextValue() (int, error) {

	var curNextVal  = math.MaxInt64
	var err error = nil

	for curNextVal, err = d.inputSets[0].Next(); err == nil;  curNextVal, err = d.inputSets[0].Next() {
		var foundVal = false
		// iterate through the rest of sets to find the diff with the first set
		// move current cursor in each set until its current value becomes equal to or greater than
		// the current value of the first set curNextVal
		// foundVal becomes true if curNextVal is found in at least one set
		for _, inputSet := range d.inputSets[1:] {

			found, err1 := readUntilValueFound(inputSet, curNextVal)
			if err1 != nil && err1 != io.EOF {
				err = err1
				break
			}
			if found { foundVal = true }
		}

		if err != nil || !foundVal { break }
	}

	return curNextVal, err
}


func readUntilValueFound(set SetReader, valToFind int) (bool, error) {

	var val int
	var err error = nil
	var found = false

	for val, err = set.Peek(); err == nil && val <= valToFind; val, err = set.Next() {
		if val == valToFind {
			found = true
			break
		}
	}
	return found, err
}