package traceerr

import (
	"errors"
	"strings"
	"testing"
)

func TestWrapPreservesErrorAndLocation(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := Wrap(sentinel)

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped error to match sentinel")
	}

	file, line, ok := Location(err)
	if !ok {
		t.Fatalf("expected wrapped error to contain location")
	}
	if !strings.HasSuffix(file, "traceerr_test.go") {
		t.Fatalf("expected test file location, got %s", file)
	}
	if line == 0 {
		t.Fatalf("expected non-zero line")
	}
}

func TestWrapNil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Fatalf("expected nil")
	}
}

func TestWrapReturnPreservesErrorAndLocation(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := wrappedReturn(sentinel)

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped error to match sentinel")
	}

	file, line, ok := Location(err)
	if !ok {
		t.Fatalf("expected wrapped error to contain location")
	}
	if !strings.HasSuffix(file, "traceerr_test.go") {
		t.Fatalf("expected test file location, got %s", file)
	}
	if line == 0 {
		t.Fatalf("expected non-zero line")
	}
}

func wrappedReturn(sentinel error) (err error) {
	defer WrapReturn(&err)

	return sentinel
}
