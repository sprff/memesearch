package contextlogger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_attrs"
)

type ContextHandler struct {
	stdout slog.Handler
	file   slog.Handler
}

var _ slog.Handler = ContextHandler{}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok && len(attrs) > 0 {
		args := make([]any, 0, 2*len(attrs))
		for _, attr := range attrs {
			args = append(args, attr.Key, attr.Value)
		}
		group := slog.Group("context", args...)
		r.AddAttrs(group)
	}

	err := h.stdout.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("can't handle stdout: %w", err)
	}

	err = h.file.Handle(ctx, r)
	if err != nil {
		return fmt.Errorf("can't handle file: %w", err)
	}

	return nil
}

func (h ContextHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.file.Enabled(ctx, l)
}

func (h ContextHandler) WithAttrs(a []slog.Attr) slog.Handler {
	return ContextHandler{
		file:   h.file.WithAttrs(a),
		stdout: h.stdout.WithAttrs(a),
	}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{
		file:   h.file.WithGroup(name),
		stdout: h.stdout.WithGroup(name),
	}
}

func NewContextHandler(fileName string, options *slog.HandlerOptions) (*ContextHandler, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}

	return &ContextHandler{
		file:   slog.NewJSONHandler(file, options),
		stdout: slog.NewJSONHandler(os.Stdout, options),
	}, nil
}

func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	var v []slog.Attr
	if va, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = va
	}

	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}
