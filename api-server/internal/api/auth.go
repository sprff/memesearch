package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"memesearch/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *api) AuthRegister(ctx context.Context, login string, password string) (models.UserID, error) {
	password = a.hashPassword(login, password)
	id, err := a.storage.CreateUser(ctx, login, password)
	if err != nil {
		if err == models.ErrUserLoginAlreadyExists {
			return "", ErrLoginExists
		}
		return models.UserID(""), fmt.Errorf("can't create: %w", err)
	}
	slog.InfoContext(ctx, "New user registered",
		"id", id,
		"login", login)
	err = a.Subscribe(ctx, id, "default", "sub")
	if err != nil {
		slog.WarnContext(ctx, "Can't subscribe to default", "err", err)
	}

	return id, nil
}

func (a *api) AuthLogin(ctx context.Context, login string, password string) (string, error) {
	password = a.hashPassword(login, password)
	user, err := a.storage.LoginUser(ctx, login, password)
	if err != nil {
		if err == models.ErrUserNotFound {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("can't login: %w", err)
	}
	token, err := a.generateToken(user)
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}
	slog.InfoContext(ctx, "Successful login",
		"login", login,
		"token", token)

	return token, nil
}

func (a *api) AuthWhoami(ctx context.Context) (models.User, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return models.User{}, ErrUnauthorized
	}

	user, err := a.storage.GetUserByID(ctx, userID)
	if err != nil {
		if err == models.ErrUserNotFound {
			return models.User{}, ErrForbidden
		}
		return models.User{}, fmt.Errorf("can't login: %w", err)
	}

	return user, nil
}

func (a *api) hashPassword(login, password string) string {
	str := fmt.Sprintf("%s:%s:%s", login, password, a.secrets.PassSalt)
	hash := sha256.Sum256([]byte(str))
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

type Claims struct {
	UserID models.UserID `json:"user_id"`
	jwt.RegisteredClaims
}

func (a *api) generateToken(u models.User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secrets.JwtCode))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type contextKey string

func (a *api) Authorize(ctx context.Context, token string) (context.Context, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.secrets.JwtCode), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, ErrInvalidToken
	}
	userId := claims.UserID
	_, err = a.storage.GetUserByID(ctx, userId)
	if err != nil {
		if err == models.ErrUserNotFound {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	return context.WithValue(ctx, contextKey("user_id"), userId), nil
}

func GetUserID(ctx context.Context) models.UserID {
	s, _ := ctx.Value(contextKey("user_id")).(models.UserID)
	return s
}
