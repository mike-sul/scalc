package scalc

import (
	"io"
)

const (
	DifOperatorId OperatorID = "DIF"
)

type DiffOperator struct {
	inputsets []SetReader
}

func NewDiffOperator(inputsets []SetReader) SetReader {
	return NewProxy(&DiffOperator{inputsets: inputsets})
}

func (d *DiffOperator) ComputeNextValue() (int, error) {
	diff_min, err := d.inputsets[0].Peek()
	if err != nil {
		return 0, err
	}

	for {
		cont := false
		for _, inputset := range d.inputsets[1:] {

			val, err := inputset.Peek()
			if err != nil && err == io.EOF {
				continue
			}

			if val < diff_min {
				inputset.Next()
				cont = true
			}

			if val == diff_min {
				diff_min, err = inputset.Next()
				if err != nil {
					return 0, err
				}
				cont = true
				break
			}

		}
		if !cont {
			break
		}
	}

	d.inputsets[0].Next()
	return diff_min, nil
}
