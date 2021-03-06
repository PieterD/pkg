package syncutil

import (
	"context"
	"sync"
)

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

func (lock *Lock) LockTry() bool {
	select {
	case lock.c <- lock.f:
		return true
	default:
		return false
	}
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

func (lock *Lock) SafeUnlock() {
	select {
	case <-lock.c:
	default:
	}
}

// TempUnlock unlocks the locker, runs f and returns what it returns, and re-locks after f has run.
func TempUnlock(locker sync.Locker, f func() error) error {
	locker.Unlock()
	defer locker.Lock()
	return f()
}

func UnlockOnce(locker sync.Locker) func() {
	once := sync.Once{}
	return func() {
		once.Do(func() {
			locker.Unlock()
		})
	}
}
