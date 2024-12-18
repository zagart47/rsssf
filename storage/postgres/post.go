package postgres

import (
	"context"
	"rsssf/entity"

	"github.com/jackc/pgx/v5"
)

type PostStorage struct {
	db Client
}

func NewPostStorage(db Client) PostStorage {
	return PostStorage{db: db}
}

func (s *PostStorage) AddPosts(ctx context.Context, posts []entity.Post) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	batch := new(pgx.Batch)
	for _, post := range posts {
		batch.Queue(`
			INSERT INTO posts (title, description, created, link)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (link) DO NOTHING
			`, post.Title, post.Content, post.PubTime, post.Link)
	}
	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostStorage) GetPosts(ctx context.Context, params int) ([]entity.Post, error) {
	limit := 10
	if params > 0 {
		limit = params
	}
	rows, err := s.db.Query(ctx, `
        SELECT id, title, description, created, link
        FROM posts
        ORDER BY created DESC
        LIMIT $1
        `, limit)
	if err != nil {
		return nil, err
	}
	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
