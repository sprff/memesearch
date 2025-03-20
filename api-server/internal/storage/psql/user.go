package psql

import (
	"context"
	"memesearch/internal/config"
	"memesearch/internal/models"

	"github.com/jmoiron/sqlx"
)

var _ models.UserRepo = &UserStore{}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(cfg config.DatabaseConfig) (UserStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return UserStore{}, err
	}
	return UserStore{db: db}, nil
}

// InsertUser implements models.UserRepo.
func (u *UserStore) InsertUser(ctx context.Context, user models.User) (models.UserID, error) {
	panic("unimplemented")
}

// GetUserByID implements models.UserRepo.
func (u *UserStore) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	panic("unimplemented")
}

// GetUserByLogin implements models.UserRepo.
func (u *UserStore) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	panic("unimplemented")
}

// UpdateUser implements models.UserRepo.
func (u *UserStore) UpdateUser(ctx context.Context, user models.User) error {
	panic("unimplemented")
}

// DeleteUser implements models.UserRepo.
func (u *UserStore) DeleteUser(ctx context.Context, id models.UserID) error {
	panic("unimplemented")
}
