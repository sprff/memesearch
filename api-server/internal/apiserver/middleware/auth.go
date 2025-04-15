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

type contextKey string

func Auth(a *api.API) func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
		return func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
				token := auth[7:]
				userID, err := a.ValidateToken(token)
				if err == nil {
					ctx := context.WithValue(r.Context(), contextKey("user_id"), string(userID))
					ctx = contextlogger.AppendCtx(ctx, slog.String("user_id", string(userID)))
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

func GetAuthUserID(ctx context.Context) string {
	s, _ := ctx.Value(contextKey("user_id")).(string)
	return s
}
