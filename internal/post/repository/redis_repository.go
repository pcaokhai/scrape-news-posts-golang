package repository

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/post"
	"github.com/redis/go-redis/v9"
)

type postsRedisRepo struct{
	redisClient *redis.Client
}

func NewPostsRedisRepo (redisClient *redis.Client) post.PostRedisRepository {
	return &postsRedisRepo{redisClient: redisClient}
}

func (p *postsRedisRepo) GetPosts(ctx context.Context, key string, out interface{}) error {
	postsBytes, err := p.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return errors.New("No key set")
	}

	if err = json.Unmarshal(postsBytes, out); err != nil {
		log.Println("Unmarshal json error")
		return errors.New("Unmarshal json error")
	}

	return nil
}

func (p *postsRedisRepo) SetPosts(ctx context.Context, key string, seconds int, posts []*models.Post) error {
	postsBytes, err := json.Marshal(posts)
	if err != nil {
		log.Println("Marshal json error")
		return errors.New("Marshal json error")
	}

	if err = p.redisClient.Set(ctx, key, postsBytes, time.Second * time.Duration(seconds)).Err(); err != nil {
		log.Printf("Set post key failed")
		return errors.New("Set post key failed")
	}

	return nil
}
