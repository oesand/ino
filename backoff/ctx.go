package backoff

import (
	"context"

	"github.com/oesand/ino/internal"
)

var ctxKey = internal.CtxKey{Key: "backoff/ctx"}

// GetContext retrieves Context from the provided context.
// It performs a safe type assertion and returns nil if the value is
// not present or not of type *Context.
func GetContext(ctx context.Context) *Context {
	data, _ := ctx.Value(ctxKey).(*Context)
	return data
}

// Context contains configuration and state for backoff attempts.
type Context struct {
	attempt, maxAttempts int
}

// Attempt returns the current retry attempt number.
func (c *Context) Attempt() int {
	return c.attempt
}

// MaxAttempts returns the maximum number of retry attempts allowed.
func (c *Context) MaxAttempts() int {
	return c.maxAttempts
}
