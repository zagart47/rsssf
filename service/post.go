package service

import (
	"context"
	"rsssf/entity"
	"rsssf/storage"
)

type PostService struct {
	postStorage storage.Storage
}

func NewPostService(postStorage storage.Storage) PostService {
	return PostService{postStorage: postStorage}
}

func (p PostService) AddPosts(ctx context.Context, posts []entity.Post) error {
	return p.postStorage.Posts.AddPosts(ctx, posts)
}

func (p PostService) GetPosts(ctx context.Context, params int) ([]entity.Post, error) {
	return p.postStorage.Posts.GetPosts(ctx, params)
}
