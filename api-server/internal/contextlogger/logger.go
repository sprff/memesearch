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

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
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

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}
