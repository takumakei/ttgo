package logger

import (
	"context"
	"log/slog"

	"github.com/gobwas/glob"
	"github.com/takumakei/ttgo/funcname"
)

type handler struct {
	handler slog.Handler
	levels  []*packageLevel
}

type packageLevel struct {
	glob  glob.Glob
	level slog.Level
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	if ctxl := FromContextOrNil(ctx); ctxl != nil {
		if ctxh := ctxl.Handler(); ctxh != h {
			return ctxh.Enabled(ctx, level)
		}
	}
	pkgname, _ := funcname.SplitCaller(4)
	for _, e := range h.levels {
		if e.glob.Match(pkgname) {
			return level >= e.level
		}
	}
	return h.handler.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if ctxl := FromContextOrNil(ctx); ctxl != nil {
		if ctxh := ctxl.Handler(); ctxh != h {
			return ctxh.Handle(ctx, r)
		}
	}
	return h.handler.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		handler: h.handler.WithAttrs(attrs),
		levels:  h.levels,
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		handler: h.handler.WithGroup(name),
		levels:  h.levels,
	}
}
