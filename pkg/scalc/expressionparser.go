package scalc

import (
	"errors"
)

func ParseExpression(exp []string, indx int) (SetReader, int, error) {
	if len(exp) < indx+4 {
		return nil, indx, errors.New("invalid expression format: too few items")
	}

	if exp[indx] != "[" {
		return nil, indx, errors.New("invalid expression format: " +
			"the specified expression doesn't start with `[` symbol")
	}

	operatorId := OperatorID(exp[indx+1])
	if !GetOperatorRegistry().IsRegistered(operatorId) {
		return nil, indx + 1, errors.New("invalid expression format: " +
			"the specified operator is not registered")
	}

	if exp[indx+2] == "]" {
		return nil, indx + 2, errors.New("invalid expression format: " +
			"none of sets are specified")
	}

	inputSets := make([]SetReader, 0, 10)
	var nextSetErr error = nil
	var curIndx int

	for curIndx = indx + 2; curIndx < len(exp) && exp[curIndx] != "]"; curIndx++ {

		var nextSet SetReader = nil

		if exp[curIndx] != "[" {
			// read the set from the specified file
			nextSet, nextSetErr = NewFileReader(exp[curIndx])
		} else {
			// got nested expression, parse it recursively
			nextSet, curIndx, nextSetErr = ParseExpression(exp, curIndx)
		}

		if nextSetErr != nil {
			break
		}
		inputSets = append(inputSets, nextSet)
	}

	if curIndx >= len(exp) {
		return nil, curIndx, errors.New("no `]` found at the end of the input expression")
	}

	if nextSetErr != nil {
		return nil, curIndx, nextSetErr
	} else {
		operator, err := GetOperatorRegistry().Create(operatorId, inputSets)
		return operator, curIndx, err
	}
}
