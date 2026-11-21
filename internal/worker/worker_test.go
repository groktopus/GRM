package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestCancellationWhileBlockedJoins(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered := make(chan struct{}, 3)
	var active, peak atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, 2, make([]int, 10), func(ctx context.Context, _ int) error {
			n := active.Add(1)
			defer active.Add(-1)
			for old := peak.Load(); n > old; old = peak.Load() {
				if peak.CompareAndSwap(old, n) {
					break
				}
			}
			entered <- struct{}{}
			<-ctx.Done()
			return ctx.Err()
		})
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-entered:
		case <-time.After(2 * time.Second):
			t.Fatal("workers did not block")
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || Classify(err) != "canceled" {
			t.Fatalf("identity lost: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown did not join")
	}
	if active.Load() != 0 || peak.Load() != 2 {
		t.Fatalf("active=%d peak=%d", active.Load(), peak.Load())
	}
}

func TestWorkerOutcomes(t *testing.T) {
	sentinel := errors.New("synthetic failure")
	err := Run(context.Background(), 2, []int{1, 2}, func(context.Context, int) error { return sentinel })
	if !errors.Is(err, sentinel) || Classify(err) != "dependency" {
		t.Fatal(err)
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	err = Run(ctx, 2, []int{1}, func(context.Context, int) error { t.Fatal("called after deadline"); return nil })
	if !errors.Is(err, context.DeadlineExceeded) || Classify(err) != "deadline" {
		t.Fatal(err)
	}
	for _, n := range []int{0, 3} {
		if Classify(Run(context.Background(), n, nil, func(context.Context, int) error { return nil })) != "input" {
			t.Fatal("bounds")
		}
	}
	var count atomic.Int32
	if err := Run(context.Background(), 2, make([]int, 10), func(context.Context, int) error { count.Add(1); return nil }); err != nil || count.Load() != 10 {
		t.Fatalf("count=%d err=%v", count.Load(), err)
	}
}

func TestDeadlineWhileBlocked(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered := make(chan struct{})
	done := make(chan error, 1)
	// The entry handshake verifies the dependency is blocked before the result.
	deadlineCtx, deadlineCancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer deadlineCancel()
	var active atomic.Int32
	go func() {
		done <- Run(deadlineCtx, 1, []int{1}, func(ctx context.Context, _ int) error {
			active.Add(1)
			defer active.Add(-1)
			close(entered)
			<-ctx.Done()
			return ctx.Err()
		})
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("not entered")
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) || Classify(err) != "deadline" || active.Load() != 0 {
			t.Fatalf("err=%v active=%d", err, active.Load())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("not joined")
	}
}

func TestBatchAdmissionAndOwnership(t *testing.T) {
	for _, workers := range []int{1, 2} {
		var called atomic.Int32
		dependency := func(context.Context, int) error { called.Add(1); return nil }
		if err := Run(context.Background(), workers, make([]int, workers+8), dependency); err != nil || called.Load() != int32(workers+8) {
			t.Fatalf("endpoint: %v calls=%d", err, called.Load())
		}
		called.Store(0)
		if err := Run(context.Background(), workers, make([]int, workers+9), dependency); !errors.Is(err, ErrInput) || called.Load() != 0 {
			t.Fatal("excess batch started effects")
		}
	}
	if err := Run(context.Background(), 1, nil, nil); !errors.Is(err, ErrInput) {
		t.Fatal("nil dependency")
	}
	if err := Run(context.Background(), 2, nil, func(context.Context, int) error { t.Fatal("empty batch called"); return nil }); err != nil {
		t.Fatal(err)
	}
	var active, peak atomic.Int32
	err := Run(context.Background(), 1, make([]int, 9), func(context.Context, int) error {
		n := active.Add(1)
		defer active.Add(-1)
		if n > peak.Load() {
			peak.Store(n)
		}
		return nil
	})
	if err != nil || peak.Load() != 1 || active.Load() != 0 {
		t.Fatal("serial ownership")
	}
}
