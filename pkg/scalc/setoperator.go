package scalc

type OperatorID string

type SetOperator interface {
	//
	ComputeNextValue() (int, error)
}
