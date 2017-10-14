package result

import "sync"

// Result is promise type interface with done channel.
type Result interface {
	Set(val interface{})
	Done() <-chan struct{}
	Value() (interface{}, error)
}

func NewResult() Result {
	res := result{
		C: make(chan struct{}),
	}
	return &res
}

type result struct {
	mu  sync.Mutex
	C   chan struct{}
	val interface{}
}

func (r *result) Set(val interface{}) {
	r.mu.Lock()
	r.val = val
	r.mu.Unlock()

	close(r.C)
}

func (r *result) Done() <-chan struct{} {
	return r.C
}

func (r *result) Wait() {
	<-r.C
}

func (r *result) Value() (interface{}, error) {
	<-r.C

	r.mu.Lock()
	err, ok := r.val.(error)
	r.mu.Unlock()

	if ok {
		return nil, err
	}
	return r.val, nil
}
