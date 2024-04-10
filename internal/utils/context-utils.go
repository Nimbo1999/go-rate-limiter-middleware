package utils

import (
	"context"
	"time"
)

func ContextWithMinimumTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 250*time.Millisecond)
}
