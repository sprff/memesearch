package models

type MediaID MemeID

type Media struct {
	ID   MediaID `json:"id"`
	Body []byte  `json:"body"`
}
