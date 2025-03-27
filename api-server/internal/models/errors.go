package models

import "errors"

// Board
var ErrBoardNotFound = errors.New("Board not found")

// Media
var ErrMediaNotFound = errors.New("Media not found")

// Meme
var ErrMemeNotFound = errors.New("Meme not found")

// User
var ErrUserNotFound = errors.New("User not found")
var ErrUserLoginAlreadyExists = errors.New("User with this login already exists")
