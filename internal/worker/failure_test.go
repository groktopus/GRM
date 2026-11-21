package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type qeError struct{}

func (*qeError) Error() string { return "qe synthetic" }
func TestQEFailureCancelsAndJoins(t *testing.T) {
	entered := make(chan struct{})
	sentinel := &qeError{}
	var active atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- Run(context.Background(), 2, []int{0, 1}, func(ctx context.Context, job int) error {
			if job == 0 {
				<-entered
				return sentinel
			}
			active.Add(1)
			defer active.Add(-1)
			close(entered)
			<-ctx.Done()
			return ctx.Err()
		})
	}()
	select {
	case err := <-done:
		var target *qeError
		if !errors.Is(err, sentinel) || !errors.As(err, &target) || active.Load() != 0 {
			t.Fatalf("err=%v active=%d", err, active.Load())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("failed sibling join")
	}
}
