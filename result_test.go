package result_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rickn42/result"
)

var (
	value   interface{}
	timeout = errors.New("timeout")
)

func TestResult_Set_and_Value(t *testing.T) {
	r := result.NewResult()
	r_err := result.NewResult()

	go func() {
		time.Sleep(10 * time.Microsecond)
		r.Set(value)
	}()

	go func() {
		time.Sleep(10 * time.Microsecond)
		r_err.Set(timeout)
	}()

	val, err := r.Value()
	if val != value || err != nil {
		t.Error("Result should be 'value'.")
	}

	val, err = r_err.Value()
	if val != nil || err != timeout {
		t.Error("Result should be 'timeout'.")
	}
}

func TestResult_Done(t *testing.T) {
	r := result.NewResult()

	select {
	case <-r.Done():
		t.Error("Result should be not done.")
	case <-time.NewTimer(10 * time.Microsecond).C:
	}

	r.Set(timeout)
	select {
	case <-r.Done():
		if _, err := r.Value(); err != timeout {
			t.Error("Result value should be 'timeout'")
		}
	case <-time.NewTimer(10 * time.Microsecond).C:
		t.Error("Result should be done.")
	}
}
