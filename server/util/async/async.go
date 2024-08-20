package async

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Go starts a goroutine with context and safely catches panics.
func Go(ctx context.Context, taskName string, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				hlog.CtxErrorf(ctx, "[async] panic in task %s: %v", taskName, r)
			}
		}()

		fn(ctx)
	}()
}