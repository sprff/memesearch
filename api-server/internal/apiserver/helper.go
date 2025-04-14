package apiserver

import (
	"fmt"
	"log/slog"
	"memesearch/internal/api"
	"memesearch/internal/models"
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

func convertMemeToModel(m Meme) (models.Meme, error) {
	dsc, err := convertMapToString(m.Description)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't convert meme: %w", err)
	}
	return models.Meme{
		ID:          models.MemeID(m.Id),
		BoardID:     models.BoardID(m.BoardId),
		Filename:    m.Filename,
		Description: dsc,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
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

func invalidInput(parametr, message string, args ...any) error {
	return &InvalidParamFormatError{ParamName: parametr, Err: fmt.Errorf(message, args...)}
}

func ptr[T any](r T) *T {
	return &r
}
