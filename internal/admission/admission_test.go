package admission

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/groktopus/GRM/internal/worker"
)

func state(o *Owner) (int, int) { o.mu.Lock(); defer o.mu.Unlock(); return o.executing, len(o.waiting) }
func until(t *testing.T, f func() bool) {
	t.Helper()
	end := time.Now().Add(3 * time.Second)
	for !f() {
		if time.Now().After(end) {
			t.Fatal("handshake timed out")
		}
		runtime.Gosched()
	}
}
func receive[T any](t *testing.T, c <-chan T) T {
	t.Helper()
	select {
	case v := <-c:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("completion timed out")
		var zero T
		return zero
	}
}
func launch(o *Owner, ctx context.Context, b Boundary, d worker.Dependency) <-chan error {
	done := make(chan error, 1)
	go func() { done <- o.Run(ctx, Query, b, []int{1}, d) }()
	return done
}
func waitState(t *testing.T, o *Owner, a, b int) {
	t.Helper()
	until(t, func() bool { x, y := state(o); return x == a && y == b })
}

func TestConcurrentCapacityAndRejectBeforeEffects(t *testing.T) {
	o := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer o.Shutdown()
	defer cancel()
	entered := make(chan struct{}, 20)
	var calls atomic.Int32
	d := func(ctx context.Context, _ int) error {
		calls.Add(1)
		entered <- struct{}{}
		<-ctx.Done()
		return ctx.Err()
	}
	first := launch(o, ctx, Ordinary, d)
	second := launch(o, ctx, Ordinary, d)
	receive(t, entered)
	receive(t, entered)
	start := make(chan struct{})
	results := make(chan error, 18)
	for i := 0; i < 18; i++ {
		go func() { <-start; results <- o.Run(ctx, Query, Ordinary, []int{1}, d) }()
	}
	close(start)
	waitState(t, o, 2, 8)
	for i := 0; i < 10; i++ {
		if err := receive(t, results); !errors.Is(err, ErrBusy) {
			t.Fatalf("excess: %v", err)
		}
	}
	if calls.Load() != 2 {
		t.Fatalf("rejected/queued calls executed: %d", calls.Load())
	}
	cancel()
	for i := 0; i < 8; i++ {
		if err := receive(t, results); !errors.Is(err, context.Canceled) {
			t.Fatalf("queued cancel: %v", err)
		}
	}
	for _, c := range []<-chan error{first, second} {
		if err := receive(t, c); !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 {
		t.Fatal("queued request invoked dependency after cancellation")
	}
	waitState(t, o, 0, 0)
}

func TestFIFOAndQueuedCancellationRelease(t *testing.T) {
	o := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer o.Shutdown()
	defer cancel()
	entered := make(chan int, 12)
	gates := make([]chan struct{}, 5)
	result := make([]<-chan error, 5)
	for i := range gates {
		gates[i] = make(chan struct{})
	}
	for i := 0; i < 5; i++ {
		id := i
		c := ctx
		var stop context.CancelFunc
		if i == 4 {
			c, stop = context.WithCancel(ctx)
			defer stop()
		}
		result[i] = launch(o, c, Ordinary, func(ctx context.Context, _ int) error {
			entered <- id
			select {
			case <-gates[id]:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
		if i < 2 {
			if got := receive(t, entered); got != i {
				t.Fatal(got)
			}
		} else {
			waitState(t, o, 2, i-1)
		}
		if i == 4 {
			stop()
			if err := receive(t, result[i]); !errors.Is(err, context.Canceled) {
				t.Fatal(err)
			}
			waitState(t, o, 2, 2)
			break
		}
	}
	// Requests 0/1 execute; 2/3 remain in original FIFO order after removing 4.
	close(gates[0])
	receive(t, result[0])
	if got := receive(t, entered); got != 2 {
		t.Fatalf("FIFO got %d", got)
	}
	close(gates[2])
	receive(t, result[2])
	if got := receive(t, entered); got != 3 {
		t.Fatalf("FIFO got %d", got)
	}
	close(gates[1])
	close(gates[3])
	receive(t, result[1])
	receive(t, result[3])
	waitState(t, o, 0, 0)
}

func TestBoundaryPermitsAcrossInvocations(t *testing.T) {
	for _, boundary := range []Boundary{Writer, Model} {
		t.Run(map[Boundary]string{Writer: "writer", Model: "model"}[boundary], func(t *testing.T) {
			o := New()
			ctx, cancel := context.WithCancel(context.Background())
			defer o.Shutdown()
			defer cancel()
			entered := make(chan struct{}, 2)
			release := make(chan struct{})
			var active, max atomic.Int32
			d := func(ctx context.Context, _ int) error {
				n := active.Add(1)
				defer active.Add(-1)
				for old := max.Load(); n > old; old = max.Load() {
					if max.CompareAndSwap(old, n) {
						break
					}
				}
				entered <- struct{}{}
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			a := launch(o, ctx, boundary, d)
			receive(t, entered)
			b := launch(o, ctx, boundary, d)
			waitState(t, o, 2, 0)
			if active.Load() != 1 {
				t.Fatal("exclusion lost")
			}
			close(release)
			receive(t, entered)
			if err := receive(t, a); err != nil {
				t.Fatal(err)
			}
			if err := receive(t, b); err != nil {
				t.Fatal(err)
			}
			if max.Load() != 1 || active.Load() != 0 {
				t.Fatal("not serialized/joined")
			}
			// Cancellation while waiting for a resource must not consume/leak its permit.
			held := make(chan struct{})
			holding := launch(o, ctx, boundary, func(ctx context.Context, _ int) error { close(held); <-ctx.Done(); return ctx.Err() })
			receive(t, held)
			short, stop := context.WithTimeout(ctx, 20*time.Millisecond)
			defer stop()
			blocked := launch(o, short, boundary, func(context.Context, int) error { t.Error("expired waiter called dependency"); return nil })
			if err := receive(t, blocked); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatal(err)
			}
			cancel()
			receive(t, holding)
			if len(o.writer) != 0 || len(o.model) != 0 {
				t.Fatal("permit leak")
			}
		})
	}
}

type domainError struct{}

func (*domainError) Error() string { return "synthetic failure" }
func TestErrorIdentityAndReuse(t *testing.T) {
	o := New()
	defer o.Shutdown()
	failure := &domainError{}
	for _, b := range []Boundary{Ordinary, Writer, Model} {
		err := o.Run(context.Background(), Prepare, b, []int{1}, func(context.Context, int) error { return failure })
		var typed *domainError
		if !errors.Is(err, failure) || !errors.As(err, &typed) || typed != failure {
			t.Fatal(err)
		}
		if err = o.Run(context.Background(), Commit, b, []int{1}, func(context.Context, int) error { return nil }); err != nil {
			t.Fatal(err)
		}
	}
	waitState(t, o, 0, 0)
}

func TestDeadlinesIncludeQueue(t *testing.T) {
	o := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer o.Shutdown()
	defer cancel()
	for _, tc := range []struct {
		op   Operation
		want time.Duration
	}{{Query, 10 * time.Second}, {Prepare, 30 * time.Second}, {Commit, 10 * time.Second}} {
		start := time.Now()
		err := o.Run(ctx, tc.op, Ordinary, []int{1}, func(c context.Context, _ int) error {
			deadline, ok := c.Deadline()
			want := tc.want
			if !ok || deadline.Before(start.Add(want)) || deadline.After(time.Now().Add(want)) {
				t.Errorf("wrong deadline %v", deadline)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	entered := make(chan struct{}, 2)
	d := func(c context.Context, _ int) error { entered <- struct{}{}; <-c.Done(); return c.Err() }
	a := launch(o, ctx, Ordinary, d)
	b := launch(o, ctx, Ordinary, d)
	receive(t, entered)
	receive(t, entered)
	short, stop := context.WithTimeout(ctx, 30*time.Millisecond)
	defer stop()
	c := launch(o, short, Ordinary, func(context.Context, int) error { t.Error("expired queued work executed"); return nil })
	waitState(t, o, 2, 1)
	if err := receive(t, c); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	waitState(t, o, 2, 0)
	cancel()
	receive(t, a)
	receive(t, b)
}

func TestShutdownCancelsAndJoinsNoncooperativeCall(t *testing.T) {
	o := New()
	entered := make(chan struct{})
	cancelled := make(chan struct{})
	release := make(chan struct{})
	var active atomic.Int32
	result := launch(o, context.Background(), Model, func(ctx context.Context, _ int) error {
		active.Add(1)
		defer active.Add(-1)
		close(entered)
		<-ctx.Done()
		close(cancelled)
		<-release
		return nil
	})
	receive(t, entered)
	stopped := make(chan struct{})
	go func() { o.Shutdown(); close(stopped) }()
	receive(t, cancelled)
	select {
	case <-stopped:
		t.Fatal("shutdown returned before join")
	default:
	}
	if err := o.Run(context.Background(), Query, Ordinary, []int{1}, func(context.Context, int) error { t.Error("closed dependency"); return nil }); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
	again := make(chan struct{})
	go func() { o.Shutdown(); close(again) }()
	close(release)
	err := receive(t, result)
	if !errors.Is(err, ErrClosed) || !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	receive(t, stopped)
	receive(t, again)
	o.Shutdown()
	if active.Load() != 0 {
		t.Fatal("active after shutdown")
	}
}

func TestShutdownRaces(t *testing.T) {
	for i := 0; i < 30; i++ {
		o := New()
		start := make(chan struct{})
		var wg sync.WaitGroup
		var active atomic.Int32
		for j := 0; j < 20; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				err := o.Run(context.Background(), Query, Model, []int{1}, func(c context.Context, _ int) error { active.Add(1); defer active.Add(-1); <-c.Done(); return c.Err() })
				if !errors.Is(err, ErrClosed) && !errors.Is(err, ErrBusy) {
					t.Errorf("race error %v", err)
				}
			}()
		}
		close(start)
		o.Shutdown()
		wg.Wait()
		o.Shutdown()
		if active.Load() != 0 {
			t.Fatal("detached dependency")
		}
		waitState(t, o, 0, 0)
	}
}

func TestInvalidAndZeroOwner(t *testing.T) {
	o := New()
	defer o.Shutdown()
	d := func(context.Context, int) error { t.Error("invalid dependency called"); return nil }
	for _, c := range []struct {
		op   Operation
		b    Boundary
		jobs []int
		d    worker.Dependency
	}{{0, Ordinary, nil, d}, {Query, Boundary(99), nil, d}, {Query, Ordinary, make([]int, 10), d}, {Query, Ordinary, nil, nil}} {
		if err := o.Run(context.Background(), c.op, c.b, c.jobs, c.d); !errors.Is(err, ErrInput) {
			t.Fatal(err)
		}
	}
	var zero Owner
	zero.Shutdown()
	if err := zero.Run(context.Background(), Query, Ordinary, []int{1}, d); !errors.Is(err, ErrClosed) {
		t.Fatal(err)
	}
}

func TestBatchStopsOnFailureAndCanceledOwnerReusable(t *testing.T) {
	o := New()
	defer o.Shutdown()
	failure := &domainError{}
	var calls atomic.Int32
	err := o.Run(context.Background(), Prepare, Writer, []int{1, 2, 3}, func(context.Context, int) error { calls.Add(1); return failure })
	if !errors.Is(err, failure) || calls.Load() != 1 {
		t.Fatalf("failure did not stop batch: %v calls=%d", err, calls.Load())
	}
	ctx, cancel := context.WithCancel(context.Background())
	entered := make(chan struct{})
	result := launch(o, ctx, Model, func(c context.Context, _ int) error { close(entered); <-c.Done(); return c.Err() })
	receive(t, entered)
	cancel()
	if err := receive(t, result); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	for _, b := range []Boundary{Writer, Model} {
		if err := o.Run(context.Background(), Query, b, []int{1}, func(context.Context, int) error { return nil }); err != nil {
			t.Fatal(err)
		}
	}
	waitState(t, o, 0, 0)
}

func TestShutdownQueuedCallsNeverStart(t *testing.T) {
	o := New()
	entered := make(chan struct{}, 2)
	var calls atomic.Int32
	d := func(c context.Context, _ int) error { calls.Add(1); entered <- struct{}{}; <-c.Done(); return c.Err() }
	a := launch(o, context.Background(), Ordinary, d)
	b := launch(o, context.Background(), Ordinary, d)
	receive(t, entered)
	receive(t, entered)
	results := make([]<-chan error, 8)
	for i := range results {
		results[i] = launch(o, context.Background(), Writer, d)
	}
	waitState(t, o, 2, 8)
	o.Shutdown()
	for _, c := range append(results, a, b) {
		if err := receive(t, c); !errors.Is(err, ErrClosed) || !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	}
	if calls.Load() != 2 {
		t.Fatal("shutdown dispatched queued dependency")
	}
	waitState(t, o, 0, 0)
}

func TestIndependentBoundaryCapacity(t *testing.T) {
	o := New()
	ctx, cancel := context.WithCancel(context.Background())
	defer o.Shutdown()
	defer cancel()
	entered := make(chan struct{}, 2)
	d := func(c context.Context, _ int) error { entered <- struct{}{}; <-c.Done(); return c.Err() }
	a := launch(o, ctx, Writer, d)
	b := launch(o, ctx, Model, d)
	receive(t, entered)
	receive(t, entered)
	cancel()
	receive(t, a)
	receive(t, b)
}
