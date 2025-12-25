package v2

import (
	"errors"
	"testing"
)

func TestJoinErr(t *testing.T) {
	t.Parallel()

	ec1 := make(chan error, 1)
	ec1 <- nil
	close(ec1)

	sentinel := errors.New("boom")
	ec2 := make(chan error, 1)
	ec2 <- sentinel
	close(ec2)

	if err := JoinErr(ec1, ec2); !errors.Is(err, sentinel) {
		t.Fatalf("err=%v, want %v", err, sentinel)
	}
}

func TestJoinErr_AllNil(t *testing.T) {
	t.Parallel()

	ec1 := make(chan error, 1)
	ec1 <- nil
	close(ec1)
	ec2 := make(chan error, 1)
	ec2 <- nil
	close(ec2)

	if err := JoinErr(ec1, ec2); err != nil {
		t.Fatalf("err=%v, want nil", err)
	}
}

func TestJoinErr_NilErrChan(t *testing.T) {
	t.Parallel()

	if err := JoinErr(nil); err == nil {
		t.Fatalf("expected error")
	} else if got, want := err.Error(), "lambda/v2: nil error channel"; got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}


