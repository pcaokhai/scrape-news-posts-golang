package post

import (
	"context"

	"github.com/pcaokhai/scraper/internal/models"
)

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	GetPostById(ctx context.Context, id int) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, post *models.Post) error
}