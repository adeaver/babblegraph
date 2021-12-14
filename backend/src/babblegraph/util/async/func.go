package async

import (
	"babblegraph/util/bglog"
	"babblegraph/util/ctx"
	"babblegraph/util/random"
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
)

type Func struct {
	fn func()
}

func WithLogContext(errs chan error, tag string, f func(c ctx.LogContext)) Func {
	return WithContext(errs, tag, func(c Context) {
		f(c)
	})
}

func WithContext(errs chan error, tag string, f func(c Context)) Func {
	contextKey := random.MustMakeRandomString(10)
	fn := func() {
		ctx := Context{
			ctx:    context.Background(),
			logger: bglog.NewLoggerForContext(tag, contextKey, 3),
		}
		defer func() {
			if x := recover(); x != nil {
				_, fn, line, _ := runtime.Caller(2)
				err := fmt.Errorf("Panic: %s: %d: %v\n%s", fn, line, x, string(debug.Stack()))
				ctx.Errorf(err.Error())
				errs <- err
			}
		}()
		f(ctx)
	}
	return Func{
		fn: fn,
	}
}

func (f Func) Start() {
	go f.fn()
}

func (f Func) Func() func() {
	return f.fn
}
