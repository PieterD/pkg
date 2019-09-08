package panic_test

import (
	"fmt"
	"testing"

	. "github.com/PieterD/pkg/panic"
)

var testError = fmt.Errorf("test error")

func TestPanic(t *testing.T) {
	err := func() (err error) {
		defer Recover(&err)
		Panic(testError)
		return nil
	}()
	if err != testError {
		t.Fatalf("expected test error, got %+v", err)
	}
}
