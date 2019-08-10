package scalc

type SetReader interface {
	// Peek returns the current `top` value of the data stream
	// without removing the value from it or if paraphrased without moving the stream cursor.
	// Returns (value, nil) on success and (math.MaxInt64, err) on failure.
	// Returns (math.MaxInt64, io.EOF) if an end of the stream is reached.
	// It's an idempotent operation until Next() is called
	Peek() (int, error)
	// Peek returns the current `top` value of the data stream
	// and removes the value from the stream (moves the stream cursor).
	// Returns (value, nil) on success and (math.MaxInt64, err) on failure.
	// Returns (math.MaxInt64, io.EOF) if an end of the stream is reached.
	Next() (int, error)
}
