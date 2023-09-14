package test

import (
	"testing"
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
