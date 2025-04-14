package psql

import (
	"context"
	"database/sql"
	"fmt"
	"memesearch/internal/config"
	"memesearch/internal/models"
	"memesearch/internal/utils"

	"github.com/jmoiron/sqlx"
)

var _ models.UserRepo = &UserStore{}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(cfg config.DatabaseConfig) (*UserStore, error) {
	db, err := connect(cfg)
	if err != nil {
		return nil, err
	}
	return &UserStore{db: db}, nil
}

// CreateUser implements models.UserRepo.
func (u *UserStore) CreateUser(ctx context.Context, login string, password string) (models.UserID, error) {
	id := utils.GenereateUUIDv7()
	var user User
	err := u.db.Get(&user, "SELECT * FROM users WHERE login=$1 LIMIT 1", login)
	if err == nil {
		return models.UserID(""), models.ErrUserLoginAlreadyExists
	}
	if err != sql.ErrNoRows {
		return models.UserID(""), fmt.Errorf("can't select: %w", err)
	}

	_, err = u.db.Exec("INSERT INTO users (id, login, password) VALUES ($1, $2, $3)", id, login, password)
	if err != nil {
		return models.UserID(""), fmt.Errorf("can't insert: %w", err)
	}
	return models.UserID(id), nil
}

// GetUserByID implements models.UserRepo.
func (u *UserStore) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	var user User
	err := u.db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("can't select: %w", err)
	}
	return convertToModelUser(user), nil
}

// LoginUser implements models.UserRepo.
func (u *UserStore) LoginUser(ctx context.Context, login string, password string) (models.User, error) {
	var user User
	err := u.db.Get(&user, "SELECT * FROM users WHERE login=$1 AND password=$2", login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("can't select: %w", err)
	}
	return convertToModelUser(user), nil
}

// UpdateUser implements models.UserRepo.
func (u *UserStore) UpdateUser(ctx context.Context, user models.User) error {
	us := convertToUser(user)
	res, err := u.db.Exec("UPDATE memes SET login = $2, passwored = $3 WHERE id=$1", us.ID, us.Login, us.Password)
	if err != nil {
		return fmt.Errorf("can't update: %w", err)
	}

	if err := zeroRows(res, models.ErrUserNotFound); err != nil {
		return err
	}
	return nil
}

// DeleteUser implements models.UserRepo.
func (u *UserStore) DeleteUser(ctx context.Context, id models.UserID) error {
	res, err := u.db.Exec("DELETE FROM memes WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("can't delete: %w", err)
	}

	if err := zeroRows(res, models.ErrUserNotFound); err != nil {
		return err
	}
	return nil
}
