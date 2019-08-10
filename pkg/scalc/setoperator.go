package scalc

type OperatorID string

type SetOperator interface {
	// Computes and returns the next value of the set that is being produced
	// by the operator implementing the given interface
	// Returns (value, nil) on success and (math.MaxInt64, err) on failure.
	// Returns (math.MaxInt64, io.EOF) if an end of the stream is reached.
	ComputeNextValue() (int, error)
}
