package middleware

import (
	"context"
	"log/slog"
	"memesearch/internal/api"
	"memesearch/internal/contextlogger"
	"net/http"
	"strings"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

func Auth(a *api.API) func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
		return func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
				token := auth[7:]

				ctx, err := a.Authorize(r.Context(), token)
				if err == nil {
					id := api.GetUserID(ctx)
					ctx = contextlogger.AppendCtx(ctx, slog.String("user_id", string(id)))
					*r = *r.WithContext(ctx)
					slog.DebugContext(ctx, "Authorized")
				} else if err == api.ErrInvalidToken {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Invalid token"))
					return nil, err
				}
			}

			res, err := f(r.Context(), w, r, request)
			return res, err
		}
	}

}
