package apiserver

import (
	"fmt"
	"log/slog"
	"memesearch/internal/api"
	"memesearch/internal/models"
	"memesearch/internal/searchranker"
	"net/http"
	"slices"
)

func NewHandler(api *api.API, middlewares []StrictMiddlewareFunc) http.Handler {
	slices.Reverse(middlewares)
	serv := NewServerImpl(api)
	serverImpl := NewStrictHandlerWithOptions(serv, middlewares, StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			slog.ErrorContext(r.Context(), "Request Error Handler", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: ErrorHandler,
	})
	handler := Handler(serverImpl)
	return handler
}

func convertMemeToServer(m models.Meme) Meme {
	dsc := convertMapToAny(m.Description)
	return Meme{
		Id:          string(m.ID),
		BoardId:     string(m.BoardID),
		Filename:    m.Filename,
		Description: dsc,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func convertScoredMemeToServer(m searchranker.ScroredMeme) ScoredMeme {
	return ScoredMeme{
		Score: m.Score,
		Meme:  convertMemeToServer(m.Meme),
	}
}

func convertMapToString(m map[string]any) (map[string]string, error) {
	dsc := map[string]string{}
	for k, v := range m {
		t, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type")
		}
		dsc[k] = t
	}
	return dsc, nil
}

func convertMapToAny(m map[string]string) map[string]any {
	dsc := map[string]any{}
	for k, v := range m {
		dsc[k] = v
	}
	return dsc
}

func convertBoardToServer(m models.Board) Board {
	return Board{
		Id:    string(m.ID),
		Owner: string(m.Owner),
		Name:  m.Name,
	}
}

func convertBoardListToServer(ms []models.Board) []Board {
	res := make([]Board, len(ms))
	for i, m := range ms {
		res[i] = convertBoardToServer(m)
	}
	return res
}

func invalidInput(parametr, message string, args ...any) error {
	return &InvalidParamFormatError{ParamName: parametr, Err: fmt.Errorf(message, args...)}
}

func ptr[T any](r T) *T {
	return &r
}
