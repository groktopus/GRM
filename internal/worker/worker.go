package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrInput identifies invalid worker configuration, before any dependency call.
var ErrInput = errors.New("invalid worker configuration")

// Dependency must return when ctx is canceled. Run joins all calls and cannot
// forcibly terminate noncooperative I/O.
type Dependency func(context.Context, int) error

// Run owns at most two workers and admits a finite batch with at most eight
// waiting jobs beyond the reserved worker slots. It is not global MCP admission.
// Use one worker for serial work; callers must separately enforce installation-wide
// writer/model exclusion. The caller must not mutate jobs until Run returns.
// There is no producer or per-job goroutine.
// On failure/cancellation it cancels all calls and joins every worker before returning.
func Run(ctx context.Context, workers int, jobs []int, dependency Dependency) error {
	if workers < 1 || workers > 2 || dependency == nil || len(jobs) > workers+8 {
		return fmt.Errorf("worker configuration: %w", ErrInput)
	}
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	var mu sync.Mutex
	next := 0
	var first error
	if workers > len(jobs) {
		workers = len(jobs)
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if child.Err() != nil || next == len(jobs) {
					mu.Unlock()
					return
				}
				job := jobs[next]
				next++
				mu.Unlock()
				if err := dependency(child, job); err != nil {
					mu.Lock()
					if first == nil {
						first = fmt.Errorf("dependency: %w", err)
					}
					mu.Unlock()
					cancel()
					return
				}
			}
		}()
	}
	wg.Wait()
	if ctx.Err() != nil {
		return fmt.Errorf("worker lifecycle: %w", ctx.Err())
	}
	return first
}

// Classify preserves the original error for errors.Is/As; it does not replace it.
func Classify(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, context.Canceled):
		return "canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "deadline"
	case errors.Is(err, ErrInput):
		return "input"
	default:
		return "dependency"
	}
}
