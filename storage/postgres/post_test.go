package postgres

import (
	"context"
	"github.com/stretchr/testify/assert"
	"rsssf/entity"
	"testing"
)

type Fields struct {
	db Client
}

var db = Fields{db: New("postgres://pgtest:pgtest@185.102.139.100:5433/pgtest?sslmode=disable")}

var news = []entity.Post{
	{
		Title:   "1st post",
		Link:    "1st url",
		Content: "1st desc",
		PubTime: 1675209600,
	},
	{
		Title:   "2nd post",
		Link:    "2nd url",
		Content: "2nd desc",
		PubTime: 1675296000,
	},
	{
		Title:   "3rd post",
		Link:    "3rd url",
		Content: "3rd desc",
		PubTime: 1675382400,
	},
}

func createTable() {
	db.db.Query(context.Background(), `
		CREATE TABLE IF NOT EXISTS posts (
    		id SERIAL PRIMARY KEY,
    		title TEXT NOT NULL,
    		description TEXT NOT NULL,
    		created BIGINT NOT NULL,
    		link varchar NOT NULL UNIQUE);
		`)

}

func dropTable() {
	db.db.Query(context.Background(), `DROP TABLE posts;`)
}

func TestPostStorage_AddPosts(t *testing.T) {
	type args struct {
		ctx   context.Context
		posts []entity.Post
	}

	tests1 := []struct {
		name    string
		fields  Fields
		args    args
		wantErr bool
	}{
		{"1st test: Провальный тест с невозможностью подключения к БД", db, args{
			posts: news, ctx: context.Background()}, true},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostStorage{
				db: tt.fields.db,
			}
			err := s.AddPosts(tt.args.ctx, tt.args.posts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	createTable()

	tests2 := []struct {
		name    string
		fields  Fields
		args    args
		wantErr bool
	}{
		{"1st test", db, args{
			posts: news, ctx: context.Background()}, false},
		{"2nd test", db, args{
			posts: nil, ctx: context.Background()}, false},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostStorage{
				db: tt.fields.db,
			}
			err := s.AddPosts(tt.args.ctx, tt.args.posts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostStorage_GetPosts(t *testing.T) {
	type fields struct {
		db Client
	}
	type args struct {
		ctx    context.Context
		params int
	}
	params := 10

	tests1 := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Post
		wantErr bool
	}{
		{"1st test",
			fields(db),
			args{
				ctx:    context.Background(),
				params: params,
			},
			news,
			false},
		{"2nd test",
			fields(db),
			args{
				ctx:    context.Background(),
				params: params,
			},
			news,
			false},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostStorage{
				db: tt.fields.db,
			}
			got, err := s.GetPosts(tt.args.ctx, tt.args.params)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, itemsEqualIgnoreID(got, tt.want), "GetPosts() got = %v, want %v", got, tt.want)
			}
		})
	}
	dropTable()
}

// itemsEqualIgnoreID сравнивает элементы структуры entity.RSSItems, игнорируя поле id.
func itemsEqualIgnoreID(got, want []entity.Post) bool {
	if len(got) != len(want) {
		return false
	}

	wantMap := make(map[entity.Post]bool)
	for _, item := range want {
		wantMap[item] = true
	}

	for _, item := range got {
		key := entity.Post{
			Title:   item.Title,
			Link:    item.Link,
			Content: item.Content,
			PubTime: item.PubTime,
		}
		if !wantMap[key] {
			return false
		}
	}

	return true
}

func TestPostStorageFail_GetPosts(t *testing.T) {
	type fields struct {
		db Client
	}
	type args struct {
		ctx    context.Context
		params int
	}
	params := 10

	dropTable()
	tests2 := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.Post
		wantErr bool
	}{
		{"1st test",
			fields(db),
			args{
				ctx:    context.Background(),
				params: params,
			},
			news,
			true},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostStorage{
				db: tt.fields.db,
			}
			_, err := s.GetPosts(tt.args.ctx, tt.args.params)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
