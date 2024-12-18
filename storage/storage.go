package storage

import (
	"context"

	"rsssf/entity"
	"rsssf/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Posts interface {
	AddPosts(context.Context, []entity.Post) error
	GetPosts(context.Context, int) ([]entity.Post, error)
}

type Storage struct {
	Posts Posts
}

func NewStorages(db *pgxpool.Pool) Storage {
	storage := postgres.NewPostStorage(db)
	return Storage{Posts: &storage}
}
