package backoff

import (
	"context"
	"errors"
	"runtime/debug"
	"time"
)

// DefaultMaxAttempts defines how many times BackOff will try the operation before giving up.
var DefaultMaxAttempts = 5

// DefaultBehaviour defines the default backoff strategy (nil by default).
// It must implement Behaviour, which determines the delay for each retry attempt.
var DefaultBehaviour Behaviour

// BackOff executes the given operation `op` with retry semantics and backoff behaviour.
//
// The function will:
//   - Execute `op` at least once.
//   - Retry up to `attempts` times (default 5) if errors occur.
//   - Wait according to the provided `Behaviour` (or `DefaultBehaviour`).
//   - Stop early if context is cancelled.
func BackOff(ctx context.Context, op func(context.Context) error, options ...Option) error {
	opts := backOffOptions{
		attempts:  DefaultMaxAttempts,
		behaviour: DefaultBehaviour,
	}

	for _, opt := range options {
		opt(&opts)
	}

	var err error

	var attempt int
	for {
		ctx := context.WithValue(ctx, ctxKey, &Context{
			attempt:     attempt,
			maxAttempts: opts.attempts,
		})

		err = op(ctx)

		if err == nil {
			break
		}

		var behaviour Behaviour

		var wrapped *Wrapped
		if errors.As(err, &wrapped) {
			behaviour = wrapped.Behaviour
		}

		if opts.attempts > 0 && attempt >= opts.attempts {
			break
		}

		if behaviour == nil {
			behaviour = opts.behaviour
		}

		if behaviour == nil {
			break
		}

		waitDuration := behaviour.Calculate(attempt)
		if waitDuration > 0 {
			select {
			case <-time.After(waitDuration):
				break
			case <-ctx.Done():
				return context.Cause(ctx)
			}
		} else {
			select {
			case <-ctx.Done():
				return context.Cause(ctx)
			default:
			}
		}

		attempt++
	}

	return err
}

// Wrap attaches a retry behaviour and stack trace to an error.
// Useful for marking errors as retryable with custom timing.
func Wrap(err error, behaviour Behaviour) error {
	var wp *Wrapped
	if errors.As(err, &wp) {
		return err
	}
	return &Wrapped{
		error:      err,
		Behaviour:  behaviour,
		StackTrace: debug.Stack(),
	}
}

// Catch helps capture Wrapped created by Wrap
func Catch(err error) (*Wrapped, error) {
	if err == nil {
		return nil, nil
	}

	var wp *Wrapped
	if errors.As(err, &wp) {
		err = wp.Unwrap()
	}

	return wp, err
}

// Wrapped wraps an error with an optional retry behaviour and stack trace.
type Wrapped struct {
	error
	StackTrace []byte
	Behaviour  Behaviour
}

// Unwrap implements the error Wrap interface
func (e *Wrapped) Unwrap() error {
	return e.error
}
