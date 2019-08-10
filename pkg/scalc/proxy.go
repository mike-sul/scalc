package scalc

type Proxy struct {
	operator SetOperator

	buffer  *int
	lastErr error
}

func NewProxy(operator SetOperator) SetReader {
	return &Proxy{operator: operator, buffer: nil, lastErr: nil}
}

func (p *Proxy) Peek() (int, error) {
	if p.buffer == nil {
		p.buffer = new(int)
		p.readValue()
	}
	return *p.buffer, p.lastErr
}

func (p *Proxy) Next() (int, error) {
	if p.buffer != nil {
		p.readValue()
	}
	return p.Peek()
}

func (p *Proxy) readValue() {
	*p.buffer, p.lastErr = p.operator.ComputeNextValue()
}
