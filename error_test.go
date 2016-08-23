package workers_test

import (
	"fmt"
	"github.com/wayt/go-workers"
	"testing"
)

func TestIsFatal(t *testing.T) {

	err := workers.Fatal("I'm a fatal error")

	if !workers.IsFatal(err) {
		t.Fatal("unrecognized fatal error")
	}
}

func TestBadIsFatal(t *testing.T) {

	err := fmt.Errorf("I'm not a fatal error")

	if workers.IsFatal(err) {
		t.Fatal("IsFatal triggered on non fatal error")
	}
}
