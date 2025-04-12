package apiserver

import (
	"log/slog"
	"memesearch/internal/api"
	"net/http"
)

func NewHandler(api *api.API) http.Handler {
	serv := NewServerImpl(api)
	serverImpl := NewStrictHandlerWithOptions(serv, nil, StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			slog.ErrorContext(r.Context(), "Request Error Handler", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: ErrorHandler,
	})
	handler := Handler(serverImpl)
	return handler
}
