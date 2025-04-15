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

func (a *API) LoginUser(ctx context.Context, login string, password string) (string, error) {
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

type Claims struct {
	UserID models.UserID `json:"user_id"`
	jwt.RegisteredClaims
}

func (a *API) generateToken(u models.User) (string, error) {
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

// ValidateToken проверяет JWT токен и возвращает UserID если токен валиден
func (a *API) ValidateToken(tokenString string) (models.UserID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.secrets.JwtCode), nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", ErrInvalidToken
}
