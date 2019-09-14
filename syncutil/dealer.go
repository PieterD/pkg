package syncutil

import (
	"context"
	"fmt"
	"sync"
)

type (
	Dealer struct {
		lock           sync.Mutex
		next           uint64
		cache          map[uint64]struct{}
		returning      bool
		returnAll      chan struct{}
		closeWhenEmpty chan struct{}
	}
	Token struct {
		dealer    *Dealer
		id        uint64
		returnAll chan struct{}
	}
)

var (
	ErrReturnAll = fmt.Errorf("return all tokens")
)

func NewDealer() *Dealer {
	closeWhenEmpty := make(chan struct{})
	close(closeWhenEmpty)
	return &Dealer{
		cache:          make(map[uint64]struct{}),
		returnAll:      make(chan struct{}),
		closeWhenEmpty: closeWhenEmpty,
	}
}

func (dealer *Dealer) Get() (*Token, error) {
	dealer.lock.Lock()
	defer dealer.lock.Unlock()

	if dealer.returning {
		return nil, fmt.Errorf("preliminary returnAll check: %w", ErrReturnAll)
	}

	if len(dealer.cache) == 0 {
		dealer.closeWhenEmpty = make(chan struct{})
	}

	id := dealer.next
	dealer.next++
	dealer.cache[id] = struct{}{}
	return &Token{
		dealer:    dealer,
		id:        id,
		returnAll: dealer.returnAll,
	}, nil
}

func (dealer *Dealer) ReturnAllTokens(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("preliminary context check: %w", ctx.Err())
	default:
	}
	dealer.lock.Lock()
	defer dealer.lock.Unlock()

	if !dealer.returning {
		dealer.returning = true
		close(dealer.returnAll)
	}

	cwe := dealer.closeWhenEmpty
	err := func() error {
		dealer.lock.Unlock()
		defer dealer.lock.Lock()
		// dealer is now temporarily unlocked

		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for returning tokens: %w", ctx.Err())
		case <-cwe:
			return nil
		}
	}()
	if err != nil {
		return err
	}

	return nil
}

// Reset
func (dealer *Dealer) Reset(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("preliminary context check: %w", ctx.Err())
	default:
	}

	var (
		returning      bool
		returnAll      <-chan struct{}
		closeWhenEmpty <-chan struct{}
	)
	func() {
		dealer.lock.Lock()
		defer dealer.lock.Unlock()

		returning = dealer.returning
		returnAll = dealer.returnAll
		closeWhenEmpty = dealer.closeWhenEmpty
	}()
	if !returning {
		return nil
	}

	select {
	default:
		return nil
	case <-returnAll:
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("waiting for tokens: %w", ctx.Err())
	case <-closeWhenEmpty:
	}

	dealer.lock.Lock()
	defer dealer.lock.Unlock()

	if dealer.returning {
		dealer.returning = false
		dealer.returnAll = make(chan struct{})
		dealer.closeWhenEmpty = make(chan struct{})
	}
	return nil
}

func (token *Token) Done() <-chan struct{} {
	return token.returnAll
}

func (token *Token) Return() {
	dealer := token.dealer
	dealer.lock.Lock()
	defer dealer.lock.Unlock()

	_, ok := dealer.cache[token.id]
	if !ok {
		return
	}
	delete(dealer.cache, token.id)
	if len(dealer.cache) == 0 {
		close(dealer.closeWhenEmpty)
	}
	return
}
