// Package admission owns the process-local interactive request lifecycle.
// Construct one Owner per installation, after the installation lock is held, and
// route every interactive invocation through it. It does not acquire that lock,
// implement administrative capacity, or establish the persistent privacy gate.
package admission

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/groktopus/GRM/internal/worker"
)

var (
	ErrBusy   = errors.New("request capacity exhausted")
	ErrClosed = errors.New("request lifecycle closed")
	ErrInput  = errors.New("unsupported request configuration")
)

type Operation uint8

const (
	Query Operation = iota + 1
	Prepare
	Commit
)

type Boundary uint8

const (
	Ordinary Boundary = iota
	Writer
	Model
)

func budget(op Operation) (time.Duration, bool) {
	switch op {
	case Query, Commit:
		return 10 * time.Second, true
	case Prepare:
		return 30 * time.Second, true
	default:
		return 0, false
	}
}

type request struct {
	ready   chan struct{}
	cancel  context.CancelCauseFunc
	running bool
}

// Owner must not be copied. Its zero value is closed; use New. No scheduler
// goroutine is needed: the caller owns its worker and release promotes FIFO work.
type Owner struct {
	mu        sync.Mutex
	closed    bool
	requests  map[*request]struct{}
	waiting   []*request
	executing int
	done      chan struct{}
	writer    chan struct{}
	model     chan struct{}
}

func New() *Owner {
	return &Owner{requests: make(map[*request]struct{}), done: make(chan struct{}),
		writer: make(chan struct{}, 1), model: make(chan struct{}, 1)}
}

// Run admits a finite internal batch as ONE interactive request. It uses one
// worker, so two executing requests cannot create more than two executing jobs.
// jobs must remain unchanged until return. All jobs share the original deadline.
// Boundary applies to each dependency call across every invocation of this owner.
// Dependencies must honor ctx and must not detach work or recursively call Run on
// this owner. Each call must describe its actual boundary; no storage or model
// adapter is implemented here. Errors retain identity; none proves commit rollback.
func (o *Owner) Run(parent context.Context, op Operation, boundary Boundary, jobs []int, dependency worker.Dependency) error {
	duration, valid := budget(op)
	if !valid || boundary > Model || dependency == nil || len(jobs) > 9 {
		return ErrInput
	}
	ctx, timeout := context.WithTimeout(parent, duration)
	defer timeout()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	r := &request{ready: make(chan struct{}), cancel: cancel}
	o.mu.Lock()
	if o.done == nil || o.closed {
		o.mu.Unlock()
		return ErrClosed
	}
	if err := ctx.Err(); err != nil {
		o.mu.Unlock()
		return lifecycleError(ctx)
	}
	if o.executing == 2 && len(o.waiting) == 8 {
		o.mu.Unlock()
		return ErrBusy
	}
	o.requests[r] = struct{}{}
	if o.executing < 2 {
		r.running = true
		o.executing++
		close(r.ready)
	} else {
		o.waiting = append(o.waiting, r)
	}
	o.mu.Unlock()
	defer o.release(r)
	select {
	case <-ctx.Done():
		return lifecycleError(ctx)
	case <-r.ready:
	}
	if ctx.Err() != nil {
		return lifecycleError(ctx)
	}
	var permit chan struct{}
	switch boundary {
	case Writer:
		permit = o.writer
	case Model:
		permit = o.model
	}
	err := worker.Run(ctx, 1, jobs, func(ctx context.Context, job int) error {
		if permit != nil {
			select {
			case permit <- struct{}{}:
			case <-ctx.Done():
				return lifecycleError(ctx)
			}
			defer func() { <-permit }()
		}
		if ctx.Err() != nil {
			return lifecycleError(ctx)
		}
		return dependency(ctx, job)
	})
	if ctx.Err() != nil {
		return lifecycleError(ctx)
	}
	return err
}

func lifecycleError(ctx context.Context) error {
	return fmt.Errorf("request lifecycle: %w", errors.Join(ctx.Err(), context.Cause(ctx)))
}

func (o *Owner) release(r *request) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.requests, r)
	if r.running {
		o.executing--
	} else {
		for i, q := range o.waiting {
			if q == r {
				o.waiting = append(o.waiting[:i], o.waiting[i+1:]...)
				break
			}
		}
	}
	if !o.closed && o.executing < 2 && len(o.waiting) > 0 {
		next := o.waiting[0]
		o.waiting = o.waiting[1:]
		next.running = true
		o.executing++
		close(next.ready)
	}
	if o.closed && len(o.requests) == 0 {
		close(o.done)
	}
}

// Shutdown atomically stops admission, cancels queued/executing requests and
// waits until their workers and dependency calls have joined. Concurrent/repeated
// calls are safe. It deliberately cannot time out while work remains detached:
// noncooperative I/O delays every Shutdown until it actually returns. Do not call
// Shutdown from a dependency running on this owner (it would wait for itself).
// Reopening/startup/privacy-gate transitions and commit reconciliation are outside
// this internal slice; construct a new owner only under the installation lifecycle.
func (o *Owner) Shutdown() {
	o.mu.Lock()
	if o.done == nil {
		o.mu.Unlock()
		return
	}
	if !o.closed {
		o.closed = true
		for r := range o.requests {
			r.cancel(ErrClosed)
		}
		if len(o.requests) == 0 {
			close(o.done)
		}
	}
	done := o.done
	o.mu.Unlock()
	<-done
}
