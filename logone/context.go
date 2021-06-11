package logone

import "context"

type key struct{}

var contextKey = &key{}

func NewContextWithLogger(parent context.Context, l *Logger) context.Context {
	return context.WithValue(parent, contextKey, l)
}

func LoggerFromContext(ctx context.Context) (*Logger, bool) {
	l, ok := ctx.Value(contextKey).(*Logger)
	return l, ok
}
