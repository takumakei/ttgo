package logger

import (
	"context"
	"log/slog"
)

type keyType struct{}

var key keyType

// Context は logger を関連づけた parent のコピーを返す.
// logger は FromContext で取り出すことができる.
func Context(parent context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(parent, key, logger)
}

// FromContext は ctx に関連づけられた *slog.Logger を返す.
// 関連づけられた *slog.Logger がなければ Root() の結果を返す.
func FromContext(ctx context.Context) *slog.Logger {
	if logger := FromContextOrNil(ctx); logger != nil {
		return logger
	}
	return Root()
}

// FromContext は ctx に関連づけられた *slog.Logger を返す.
// 関連づけられた *slog.Logger がなければ nil を返す.
func FromContextOrNil(ctx context.Context) *slog.Logger {
	if v := ctx.Value(key); v != nil {
		return v.(*slog.Logger)
	}
	return nil
}
