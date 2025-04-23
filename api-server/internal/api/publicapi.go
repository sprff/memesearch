package api

import (
	"context"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/searchranker"
	"memesearch/internal/storage"
)

type API struct {
	api *api
}

func New(s storage.Storage, secrets config.SecretConfig, ranker searchranker.Ranker) *API {
	return &API{&api{
		storage: s,
		secrets: secrets,
		ranker:  ranker,
	}}
}

func (a *API) CreateBoard(ctx context.Context, name string) (models.Board, error) {
	if GetUserID(ctx) == "" {
		return models.Board{}, ErrUnauthorized
	}
	return a.api.CreateBoard(ctx, name)
}

func (a *API) GetBoardByID(ctx context.Context, id models.BoardID) (models.Board, error) {
	if err := a.aclGetBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}

	return a.api.GetBoardByID(ctx, id)
}

func (a *API) UpdateBoard(ctx context.Context, id models.BoardID, name *string, owner *models.UserID) (models.Board, error) {
	if owner != nil {
		if err := a.validateBoard(ctx, id, "new owner"); err != nil {
			return models.Board{}, err
		}
	}

	if err := a.aclUpdateBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.UpdateBoard(ctx, id, name, owner)
}

func (a *API) DeleteBoard(ctx context.Context, id models.BoardID) (models.Board, error) {
	if err := a.aclDeleteBoard(ctx, id); err != nil {
		return models.Board{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.DeleteBoard(ctx, id)
}

func (a *API) ListBoards(ctx context.Context, offset, limit int, sortBy string) ([]models.Board, error) {
	if GetUserID(ctx) == "" {
		return nil, ErrUnauthorized
	}
	return a.api.ListBoards(ctx, offset, limit, sortBy)
}

func (a *API) GetMedia(ctx context.Context, id models.MediaID) (models.Media, error) {
	if err := a.aclGetMedia(ctx, id); err != nil {
		return models.Media{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.GetMedia(ctx, id)
}

func (a *API) SetMedia(ctx context.Context, media models.Media) error {
	if err := a.aclUpdateMedia(ctx, media.ID); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}
	return a.api.SetMedia(ctx, media)
}

func (a *API) CreateMeme(ctx context.Context, board models.BoardID, filename string, dsc map[string]string) (models.Meme, error) {
	if err := a.validateBoard(ctx, board, "meme's board"); err != nil {
		return models.Meme{}, err
	}
	if err := a.aclPostMeme(ctx, board); err != nil {
		return models.Meme{}, fmt.Errorf("acl failed: %w", err)
	}

	return a.api.CreateMeme(ctx, board, filename, dsc)
}

func (a *API) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {
	if err := a.aclGetMeme(ctx, id); err != nil {
		return models.Meme{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.GetMemeByID(ctx, id)
}

func (a *API) UpdateMeme(ctx context.Context, id models.MemeID, board *models.BoardID, filename *string, dsc *map[string]string) (models.Meme, error) {
	if board != nil {
		if err := a.validateBoard(ctx, *board, "meme's board"); err != nil {
			return models.Meme{}, err
		}
	}

	if err := a.aclUpdateMeme(ctx, id); err != nil {
		return models.Meme{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.UpdateMeme(ctx, id, board, filename, dsc)
}

func (a *API) DeleteMeme(ctx context.Context, id models.MemeID) error {
	if err := a.aclDeleteMeme(ctx, id); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}
	return a.api.DeleteMeme(ctx, id)
}

func (a *API) ListMemes(ctx context.Context, offset, limit int, sortBy string) ([]models.Meme, error) {
	return a.api.ListMemes(ctx, offset, limit, sortBy)
}

func (a *API) Unsubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	if err := a.aclUnsubscribe(ctx, user, board, role); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}
	return a.api.Unsubscribe(ctx, user, board, role)
}

func (a *API) Subscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	if _, err := a.GetBoardByID(ctx, board); err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}
	if err := a.aclSubscribe(ctx, user, board, role); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}
	return a.api.Subscribe(ctx, user, board, role)
}

func (a *API) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	if err := a.aclGetUser(ctx, id); err != nil {
		return models.User{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.GetUserByID(ctx, id)
}

func (a *API) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if err := a.aclUpdateUser(ctx, user.ID); err != nil {
		return models.User{}, fmt.Errorf("acl failed: %w", err)
	}
	return a.api.UpdateUser(ctx, user)
}

func (a *API) DeleteUser(ctx context.Context, id models.UserID) error {
	if err := a.aclDeleteUser(ctx, id); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}
	return a.api.DeleteUser(ctx, id)
}

func (a *API) AuthRegister(ctx context.Context, login string, password string) (models.UserID, error) {
	return a.api.AuthRegister(ctx, login, password)
}

func (a *API) AuthLogin(ctx context.Context, login string, password string) (string, error) {
	return a.api.AuthLogin(ctx, login, password)
}

func (a *API) AuthWhoami(ctx context.Context) (models.User, error) {
	return a.api.AuthWhoami(ctx)
}

func (a *API) Authorize(ctx context.Context, token string) (context.Context, error) {
	return a.api.Authorize(ctx, token)
}

func (a *API) Search(ctx context.Context, req map[string]string, offset, limit int) ([]searchranker.ScroredMeme, error) {
	if len(req) == 0 {
		return nil, ErrInvalid{"reqest", "request shouldn't be empty"}
	}
	return a.api.Search(ctx, req, offset, limit)
}
