package post

import (
	"context"

	"github.com/pcaokhai/scraper/internal/models"
)

type PostRedisRepository interface {
	GetPosts(ctx context.Context, key string, out interface{}) error
	SetPosts(ctx context.Context, key string, seconds int, posts []*models.Post) error
}
