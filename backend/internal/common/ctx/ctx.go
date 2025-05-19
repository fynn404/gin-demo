package ctx

import (
	"context"
	"time"
)

// DefaultTimeout 默认超时时间
const DefaultTimeout = 5 * time.Second

// Background 返回一个带有默认超时的 context
func Background() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// BackgroundWithTimeout 返回一个带有指定超时的 context
func BackgroundWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
