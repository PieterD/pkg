package syncutil

import (
	"context"
	"errors"
	"testing"
	"time"
)

func isAlive(t *testing.T, toks ...*Token) {
	t.Helper()
	for _, tok := range toks {
		select {
		case <-tok.Done():
			t.Fatalf("expected alive, but tok was done: %#v", tok)
		default:
		}
	}
}

func isDealing(t *testing.T, d *Dealer) {
	t.Helper()
	select {
	case <-d.Dealing():
	default:
		t.Fatalf("dealer should be dealing")
	}
}

func isNotDealing(t *testing.T, d *Dealer) {
	t.Helper()
	select {
	case <-d.Dealing():
		t.Fatalf("dealer should not be dealing")
	default:
	}
}

func TestDealer_Immediate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	d := NewDealer()
	if err := d.Reset(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled error, got: %v", err)
	}
	if err := d.ReturnAllTokens(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Canceled error, got: %v", err)
	}
	if err := d.ReturnAllTokens(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.Reset(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.Reset(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.ReturnAllTokens(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.Reset(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.ReturnAllTokens(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if err := d.Reset(nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDealer(t *testing.T) {
	d := NewDealer()
	isDealing(t, d)
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
	t1, err := d.Get()
	if err != nil {
		t.Fatalf("error getting token 1: %v", err)
	}
	t2, err := d.Get()
	if err != nil {
		t.Fatalf("error getting token 2: %v", err)
	}
	isAlive(t, t1, t2)
	go func() {
		defer t1.Return()
		<-t1.Done()
	}()
	go func() {
		defer t2.Return()
		<-t2.Done()
		t2.Return()
		t2.Return()
		t2.Return()
	}()
	isDealing(t, d)
	if err := d.ReturnAllTokens(ctx); err != nil {
		t.Fatalf("failed to return all tokens: %v", err)
	}
	isNotDealing(t, d)
	_, err = d.Get()
	if !errors.Is(err, ErrReturnAll) {
		t.Fatalf("expected ReturnAll error, got: %v", err)
	}
	isNotDealing(t, d)
	if err := d.Reset(ctx); err != nil {
		t.Fatalf("failed to reset dealer: %v", err)
	}
	isDealing(t, d)
	if err := d.Reset(ctx); err != nil {
		t.Fatalf("failed to reset dealer: %v", err)
	}
	isDealing(t, d)
	_, err = d.Get()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	isDealing(t, d)
	if err := d.Reset(ctx); err != nil {
		t.Fatalf("failed to reset dealer: %v", err)
	}
	isDealing(t, d)
	if err := d.ReturnAllTokens(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded error, got: %v", err)
	}
	isNotDealing(t, d)
	ctx, _ = context.WithTimeout(context.Background(), time.Millisecond*100)
	if err := d.ReturnAllTokens(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded error, got: %v", err)
	}
	isNotDealing(t, d)
}
