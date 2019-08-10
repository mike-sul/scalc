package scalc

import (
	"io"
	"math"
)

type SetVal struct {
	val int
	err error
}

type ChannelProxy struct {
	operator SetOperator

	pipe chan SetVal
	lastErr error
	buffer  *int
}

func NewChannelProxy(operator SetOperator) SetReader {
	cp :=  ChannelProxy{operator: operator, pipe: make(chan SetVal), lastErr: nil, buffer:nil}
	go cp.readValue()
	return &cp
}

func (p *ChannelProxy) Peek() (int, error) {
	if p.buffer == nil {
		p.buffer = new(int)
		p.readValueFromChannel()
	}
	return *p.buffer, p.lastErr
}

func (p *ChannelProxy) Next() (int, error) {
	if p.buffer != nil {
		p.readValueFromChannel()
	}
	return p.Peek()
}

func (p *ChannelProxy) readValue() {

	var val int
	var err error
	for val, err = p.operator.ComputeNextValue(); err == nil; val, err = p.operator.ComputeNextValue() {
		p.pipe <- SetVal{val: val, err: err}
	}
	p.pipe <- SetVal{val: val, err: err}
	close(p.pipe)
}


func (p *ChannelProxy) readValueFromChannel() {
	val, ok := <-p.pipe
	if ok {
		p.lastErr = val.err
		*p.buffer = val.val
	} else {
		p.lastErr = io.EOF
		*p.buffer = math.MaxInt64
	}
}