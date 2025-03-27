package models

type UserID string

type User struct {
	ID       UserID `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
