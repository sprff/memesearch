package models

type BoardID string

type Board struct {
	ID    BoardID `json:"id"    db:"id"`
	Owner UserID  `json:"owner" db:"owner_id"`
	Name  string  `json:"name"  db:"name"`
}
