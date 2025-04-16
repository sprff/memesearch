package models

import "context"

type Subsciption struct {
	BoardID BoardID
	UserID  UserID
	Role    string
}

type SubsciptionRepo interface {
	Subscribe(ctx context.Context, user UserID, board BoardID, role string) error
	Unsubscribe(ctx context.Context, user UserID, board BoardID, role string) error
}
