package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

// ---BOARD---

func (a *API) aclGetBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	return nil
}

func (a *API) aclUpdateBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	board, err := a.api.GetBoardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}
	if board.Owner != userID {
		return ErrForbidden
	}
	return nil
}

func (a *API) aclDeleteBoard(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	board, err := a.api.GetBoardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}
	if board.Owner != userID {
		return ErrForbidden
	}

	return nil
}

// ---BOARD---
// ---MEME---

func (a *API) aclPostMeme(ctx context.Context, id models.BoardID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	err := a.aclUpdateBoard(ctx, id)
	if err != nil {
		return fmt.Errorf("acl update board failed: %w", err)
	}

	return nil
}

func (a *API) aclGetMeme(ctx context.Context, id models.MemeID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	meme, err := a.api.GetMemeByID(ctx, id)
	if err != nil {
		if err == models.ErrMemeNotFound {
			return ErrMemeNotFound
		}
		return fmt.Errorf("can't get meme: %w", err)
	}

	board, err := a.api.GetBoardByID(ctx, meme.BoardID)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}

	err = a.aclGetBoard(ctx, board.ID)
	if err != nil {
		return fmt.Errorf("acl get board failed: %w", err)
	}

	return nil
}

func (a *API) aclUpdateMeme(ctx context.Context, id models.MemeID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	meme, err := a.api.GetMemeByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get meme: %w", err)
	}

	board, err := a.api.GetBoardByID(ctx, meme.BoardID)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}

	err = a.aclUpdateBoard(ctx, board.ID)
	if err != nil {
		return fmt.Errorf("acl update board failed: %w", err)
	}

	return nil
}

func (a *API) aclDeleteMeme(ctx context.Context, id models.MemeID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}

	meme, err := a.api.GetMemeByID(ctx, id)
	if err != nil {
		return fmt.Errorf("can't get meme: %w", err)
	}

	board, err := a.api.GetBoardByID(ctx, meme.BoardID)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}

	err = a.aclDeleteBoard(ctx, board.ID)
	if err != nil {
		return fmt.Errorf("acl update board failed: %w", err)
	}

	return nil
}

// ---MEME---
// ---MEDIA---

func (a *API) aclGetMedia(ctx context.Context, id models.MediaID) error {
	err := a.aclGetMeme(ctx, models.MemeID(id))
	if err != nil {
		return fmt.Errorf("acl update meme failed: %w", err)
	}
	return nil
}

func (a *API) aclUpdateMedia(ctx context.Context, id models.MediaID) error {
	err := a.aclUpdateMeme(ctx, models.MemeID(id))
	if err != nil {
		return fmt.Errorf("acl update meme failed: %w", err)
	}
	return nil
}

// ---MEDIA---
// ---USER---

func (a *API) aclGetUser(ctx context.Context, id models.UserID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	if userID != id {
		return ErrForbidden
	}
	return nil

}

func (a *API) aclUpdateUser(ctx context.Context, id models.UserID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	if userID != id {
		return ErrForbidden
	}
	return nil

}

func (a *API) aclDeleteUser(ctx context.Context, id models.UserID) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	if userID != id {
		return ErrForbidden
	}
	return nil

}

// ---USER---
// ---SUBS---

func (a *API) aclSubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	if userID != user {
		return ErrForbidden
	}
	// TODO  fail if board is private
	return nil

}

func (a *API) aclUnsubscribe(ctx context.Context, user models.UserID, board models.BoardID, role string) error {
	userID := GetUserID(ctx)
	if userID == "" {
		return ErrUnauthorized
	}
	if userID == user { // Self unsubscribe
		return nil
	}

	b, err := a.api.GetBoardByID(ctx, board)
	if err != nil {
		return fmt.Errorf("can't get board: %w", err)
	}

	if userID == b.Owner { // Owner can delete subscribers
		return nil
	}

	return ErrForbidden
}

// ---SUBS---
