package scalc

import (
	"errors"
	"sync"
)

type OperatorIdToCtorMap map[OperatorID]func([]SetReader) SetReader

type OperatorRegistry struct {
	factories OperatorIdToCtorMap
}

var OperatorRegistrySingleton *OperatorRegistry
var once sync.Once

func GetOperatorRegistry() *OperatorRegistry {
	once.Do(func() {
		OperatorRegistrySingleton = &OperatorRegistry{
			factories: OperatorIdToCtorMap{
				UnionOperatorId: NewUnionOperator,
				InterOperatorId: NewInterOperator,
				DifOperatorId:   NewDiffOperator,
			},
		}
	})
	return OperatorRegistrySingleton
}

func (or *OperatorRegistry) IsRegistered(id OperatorID) bool {
	_, found := or.factories[id]
	return found
}

func (or *OperatorRegistry) Create(id OperatorID, inputsets []SetReader) (SetReader, error) {
	ctor, found := or.factories[id]
	if !found {
		return nil, errors.New("the requested operator is not registered: " + string(id))
	}
	//TODO: an operator ctor/new might fail, consider returning (SetReader, error)
	return ctor(inputsets), nil
}
