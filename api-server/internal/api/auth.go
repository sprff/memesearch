package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
)

func (a *API) LoginUser(ctx context.Context, login string, password string) (string, error) {
	password = a.hashPassword(login, password)
	user, err := a.storage.LoginUser(ctx, login, password)
	if err != nil {
		if err == models.ErrUserNotFound {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("can't login: %w", err)
	}
	slog.InfoContext(ctx, "Successful login",
		"login", login)
	return string(user.ID), nil
}

func (a *API) Whoami(ctx context.Context, token string) (models.User, error) {
	user, err := a.storage.GetUserByID(ctx, models.UserID(token))
	if err != nil {
		if err == models.ErrUserNotFound {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("can't login: %w", err)
	}

	return user, nil
}

func (a *API) hashPassword(login, password string) string {
	str := fmt.Sprintf("%s:%s:%s", login, password, "SALT") // TODO use secret.Salt
	hash := sha256.Sum256([]byte(str))
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

func (a *API) ValidateToken(token string) (models.UserID, error) {
	//TODO jwt
	return models.UserID("019634b6d74e7426a26b5f10e3b90f5f"), nil
}
