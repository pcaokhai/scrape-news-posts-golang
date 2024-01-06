package post

import (
	"context"

	"github.com/pcaokhai/scraper/internal/models"
)

type PostUseCase interface {
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	UpdatePost(ctx context.Context, id int, title string) (*models.Post, error)
	DeletePost(ctx context.Context, id int) error
}