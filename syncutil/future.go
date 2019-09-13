package syncutil

import "fmt"

type (
	Future interface {
		Done() <-chan struct{}
		Err() error
	}
	Finisher interface {
		Future
		Finish(err error)
	}
	future struct {
		done chan struct{}
		err  error
	}
)

var (
	ErrInProgress = fmt.Errorf("future still in progress")
)

func NewFuture() Finisher {
	return &future{
		done: make(chan struct{}),
	}
}

func (f *future) Done() <-chan struct{} {
	return f.done
}

func (f *future) Err() error {
	select {
	case <-f.done:
		return f.err
	default:
		return ErrInProgress
	}
}

func (f *future) Finish(err error) {
	select {
	case <-f.done:
		return
	default:
		f.err = err
		close(f.done)
	}
}
