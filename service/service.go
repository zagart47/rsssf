package service

import (
	"context"
	"rsssf/entity"
	"rsssf/storage"
)

type Poster interface {
	AddPosts(ctx context.Context, posts []entity.Post) error
	GetPosts(ctx context.Context, params int) ([]entity.Post, error)
}

type Services struct {
	Poster Poster
}

func NewServices(storage storage.Storage) Services {
	postServices := NewPostService(storage)
	return Services{
		Poster: postServices,
	}
}
