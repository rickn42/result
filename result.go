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
	once sync.Once
	C    chan struct{}
	mu   sync.Mutex
	val  interface{}
	err  error
}

func (r *result) closeCh() {
	r.once.Do(func() {
		close(r.C)
	})
}

func (r *result) Set(val interface{}) {
	if err, ok := val.(error); ok {
		r.mu.Lock()
		r.val = nil
		r.err = err
		r.mu.Unlock()
	} else {
		r.mu.Lock()
		r.val = val
		r.err = nil
		r.mu.Unlock()
	}
	r.closeCh()
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
	val, err = r.val, r.err
	r.mu.Unlock()

	return val, err
}

func (r *result) Err() (err error) {
	<-r.C

	r.mu.Lock()
	err = r.err
	r.mu.Unlock()

	return err
}
