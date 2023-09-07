package test

import (
	"database/sql/driver"
	"testing"
	"time"
)

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err: %e", err)
	}
}

func Equal[T comparable](t *testing.T, x, y T) {
	if x != y {
		t.Fatalf("not equal: %v, %v", x, y)
	}
}

// --- For sqlmock ---

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
