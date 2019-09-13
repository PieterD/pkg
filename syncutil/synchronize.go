package syncutil

import "context"

type (
	Lock struct {
		c chan func()
		f func()
	}
)

//TODO: expand into Semaphore
func NewLock() *Lock {
	c := make(chan func(), 1)
	var f func()
	f = func() { c <- f }
	return &Lock{
		c: c,
		f: f,
	}
}

func (lock *Lock) Lock() {
	lock.c <- lock.f
}

func (lock *Lock) LockCtx(ctx context.Context) error {
	if ctx == nil {
		lock.Lock()
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case lock.c <- lock.f:
		return nil
	}
}

func (lock *Lock) Unlock() {
	<-lock.c
}
