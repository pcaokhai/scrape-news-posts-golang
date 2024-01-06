package repository

import (
	"context"

	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/post"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// GetPosts implements post.PostRedisRepository.
func (*postRepository) GetPosts(ctx context.Context, key string) ([]*models.Post, error) {
	panic("unimplemented")
}

// SetPosts implements post.PostRedisRepository.
func (*postRepository) SetPosts(ctx context.Context, key string, seconds int, news []*models.Post) error {
	panic("unimplemented")
}

func NewPostRepository(db *gorm.DB) post.PostRepository {
	return &postRepository{db: db}
}

func (pr *postRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	var posts []*models.Post

	err := pr.db.WithContext(ctx).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *postRepository) GetPostById(ctx context.Context, id int) (*models.Post, error) {
	var post *models.Post

	err := pr.db.WithContext(ctx).Where(&models.Post{Id: id}).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *postRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	err := pr.db.WithContext(ctx).Save(&post).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *postRepository) DeletePost(ctx context.Context, post *models.Post) error {
	err := pr.db.WithContext(ctx).Delete(post).Error
	if err != nil {
		return err
	}
	return nil
}
