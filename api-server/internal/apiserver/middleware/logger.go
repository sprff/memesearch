package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"memesearch/internal/contextlogger"
	"memesearch/internal/utils"
	"net/http"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

type Middleware = strictnethttp.StrictHTTPMiddlewareFunc

func Logger() func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
		return func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = utils.GenereateUUIDv7()
			}
			w.Header().Add("X-Request-ID", id)

			ctx := contextlogger.AppendCtx(r.Context(), slog.String("request_id", id))
			slog.InfoContext(ctx, fmt.Sprintf("%s %s start", r.Method, r.RequestURI))
			*r = *r.WithContext(ctx)
			res, err := f(ctx, w, r, request)
			if err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("%s %s failed", r.Method, r.RequestURI), "err", err)
			} else {
				slog.InfoContext(ctx, fmt.Sprintf("%s %s finished", r.Method, r.RequestURI))
			}

			return res, err
		}
	}
}
