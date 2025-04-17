package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) validateUser(ctx context.Context, id models.UserID, param string) error {
	if _, err := a.api.GetUserByID(ctx, id); err != nil {
		if err == ErrUserNotFound {
			return ErrInvalid{Param: param, Reason: "user don't exist"}
		}
		return fmt.Errorf("can't get user: %w", err)
	}
	return nil
}

func (a *API) validateBoard(ctx context.Context, id models.BoardID, param string) error {
	if _, err := a.api.GetBoardByID(ctx, id); err != nil {
		if err == ErrBoardNotFound {
			return ErrInvalid{Param: param, Reason: "board don't exist"}
		}
		return fmt.Errorf("can't get board: %w", err)
	}
	return nil
}

func (a *API) validateMeme(ctx context.Context, id models.MemeID, param string) error {
	if _, err := a.api.GetMemeByID(ctx, id); err != nil {
		if err == ErrMemeNotFound {
			return ErrInvalid{Param: param, Reason: "meme don't exist"}
		}
		return fmt.Errorf("can't get meme: %w", err)
	}
	return nil
}
