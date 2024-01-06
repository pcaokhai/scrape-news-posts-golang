package usecase

import (
	"context"
	"log"

	"github.com/pcaokhai/scraper/pkg/utils"
	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/post"
)

type postUseCase struct {
	postRepo 	post.PostRepository
	redisRepo 	post.PostRedisRepository
}

func NewPostUseCase(postRepo post.PostRepository, redisRepo post.PostRedisRepository) post.PostUseCase {
	return &postUseCase{postRepo: postRepo, redisRepo: redisRepo}
}

func (pu *postUseCase) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	var cachedPost []*models.Post

	// check if posts are already cached
	if err := pu.redisRepo.GetPosts(ctx, utils.PostsKey, &cachedPost); err == nil {
		log.Printf("data returned from cache")
		return cachedPost, nil
	}

	// fetch the posts if no cache
	posts, err := pu.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	// cache the posts for later query
	if err := pu.redisRepo.SetPosts(ctx, utils.PostsKey, utils.CacheDuration, posts); err != nil {
		log.Printf("error caching posts")
	}

	log.Printf("posts cached")
	return posts, nil
}

func (pu *postUseCase) UpdatePost(ctx context.Context, id int, title string) (*models.Post, error) {
	post, err := pu.postRepo.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	post.Title = title

	// update post
	err = pu.postRepo.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	// re-fetch all posts to re-cache
	posts, err := pu.postRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	// re-cache posts
	if err := pu.redisRepo.SetPosts(ctx, utils.PostsKey, utils.CacheDuration, posts); err != nil {
		log.Printf("error cache posts")
	}

	return post, nil
}

func (pu *postUseCase) DeletePost(ctx context.Context, id int) error{
	post, err := pu.postRepo.GetPostById(ctx, id)
	if err != nil {
		return err
	}

	err = pu.postRepo.DeletePost(ctx, post)
	if err != nil {
		return err
	}

	// after delete a post, need to re-fetch all posts to re-cache
	posts, err := pu.postRepo.GetAllPosts(ctx)
	if err != nil {
		return err
	}

	// re-cache posts
	if err := pu.redisRepo.SetPosts(ctx, utils.PostsKey, utils.CacheDuration, posts); err != nil {
		log.Printf("error cache posts")
	}

	return nil
}