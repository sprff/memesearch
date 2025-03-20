package storage

import "memesearch/internal/models"

type Storage struct {
	models.BoardRepo
	models.MemeRepo
	models.MediaRepo
	models.UserRepo
}
