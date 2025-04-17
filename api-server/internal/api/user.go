package api

import (
	"context"
	"fmt"
	"memesearch/internal/models"
)

func (a *API) GetUser(ctx context.Context, id models.UserID) (models.User, error) {
	if err := a.aclGetUser(ctx, id); err != nil {
		return models.User{}, fmt.Errorf("acl failed: %w", err)
	}

	user, err := a.storage.GetUserByID(ctx, id)
	if err != nil {
		if err == models.ErrUserNotFound {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("can't get: %w", err)
	}
	return user, nil
}

func (a *API) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	if err := a.aclUpdateUser(ctx, user.ID); err != nil {
		return models.User{}, fmt.Errorf("acl failed: %w", err)
	}

	err := a.storage.UpdateUser(ctx, user)
	if err != nil {
		if err == models.ErrUserNotFound {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("can't login: %w", err)
	}

	user, err = a.storage.GetUserByID(ctx, user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("can't login: %w", err)
	}
	return user, nil
}

func (a *API) DeleteUser(ctx context.Context, id models.UserID) error {
	if err := a.aclDeleteUser(ctx, id); err != nil {
		return fmt.Errorf("acl failed: %w", err)
	}

	err := a.storage.DeleteUser(ctx, id)
	if err != nil {
		if err == models.ErrUserNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("can't login: %w", err)
	}
	return nil
}
