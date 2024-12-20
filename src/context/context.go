package context

import (
	"context"
	"log/slog"
)

const callersKey = "callers"

func AppendCallStack(ctx context.Context, callstack string) context.Context {
	_callstack, _ := ctx.Value(callersKey).(string)
	if _callstack != "" {
		ctx = context.WithValue(ctx, callersKey, _callstack+"->"+callstack)
	} else {
		ctx = context.WithValue(ctx, callersKey, callstack)
	}
	return ctx
}

func AsAttributes(ctx context.Context) []any {
	_callstack, _ := ctx.Value(callersKey).(string)
	callstack := slog.String(callersKey, _callstack)

	return []any{
		callstack,
	}
}
