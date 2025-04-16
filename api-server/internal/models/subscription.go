package models

import "context"

type Subsciption struct {
	BoardID BoardID
	UserID  UserID
	Role    string
}

type SubsciptionRepo interface {
	Subscribe(ctx context.Context, sub Subsciption) error
	Unsubscribe(ctx context.Context, sub Subsciption) error
}
