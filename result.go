package result

import "sync"

// Result is promise type interface with done channel.
type Result interface {
	Set(val interface{})
	Done() <-chan struct{}
	IsDone() bool
	Wait()
	Value() (interface{}, error)
	Err() error
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

func (r *result) IsDone() (ok bool) {
	select {
	case <-r.C:
		ok = true
	default:
	}

	return ok
}

func (r *result) Wait() {
	<-r.C
}

func (r *result) Value() (val interface{}, err error) {
	<-r.C

	r.mu.Lock()
	err, ok := r.val.(error)
	r.mu.Unlock()

	if ok {
		return nil, err
	}

	val = r.val
	return val, nil
}

func (r *result) Err() (err error) {
	<-r.C

	r.mu.Lock()
	err, ok := r.val.(error)
	r.mu.Unlock()

	if ok {
		return err
	}

	return nil
}
