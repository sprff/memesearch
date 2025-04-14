package models

import (
	"context"
)

type UserID string

type User struct {
	ID       UserID `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, login, password string) (UserID, error)
	GetUserByID(ctx context.Context, id UserID) (User, error)
	LoginUser(ctx context.Context, login, password string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id UserID) error
}
