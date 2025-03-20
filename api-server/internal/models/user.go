package models

import (
	"context"
	"errors"
)

type UserID string

type User struct {
	ID       UserID `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRepo interface {
	InsertUser(ctx context.Context, user User) (UserID, error)
	GetUserByID(ctx context.Context, id UserID) (User, error)
	GetUserByLogin(ctx context.Context, login string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id UserID) error
}

var ErrUserNotFound = errors.New("User not found")
var ErrUserLoginAlreadyExists = errors.New("User with this login already exists")
