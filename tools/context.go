package tools

import (
	"context"
	"time"
)

func ContextClose(fn context.CancelFunc) {
	if fn == nil {
		return
	}
	fn()
}

func ContextBackground() context.Context {
	return context.Background()
}

func ContextToDo() context.Context {
	return context.TODO()
}

func NewContextTimeOut(ctx context.Context, timeOut time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeOut)
}

func NewContextCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

func NewContextValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(context.Background(), key, value)
}
