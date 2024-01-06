package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/post"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	postsKey = "posts"
	cacheDuration = 10				// cache the posts for 30s (for demo purpose)
)

func GetCurrentDate() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%v-%v-%v", year, int(month), day)
}

func FetchData(ctx context.Context ,db *gorm.DB, rdb *redis.Client, redisRepo post.PostRedisRepository)  ([]*models.Post, error) {
	var cachedPost []*models.Post

	// check if posts are already cached
	if err := redisRepo.GetPosts(ctx, PostsKey, &cachedPost); err == nil {
		log.Printf("data returned from cache")
		return cachedPost, nil
	}

	// cache the posts for later query
	var posts []*models.Post
	err := db.WithContext(ctx).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	// cache the posts for later query
	if err := redisRepo.SetPosts(ctx, PostsKey, CacheDuration, posts); err != nil {
		log.Printf("error caching posts")
	}
	
	log.Printf("posts cached")
	return posts, nil
}