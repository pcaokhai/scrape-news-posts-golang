package redis

import (
	"fmt"

	"github.com/pcaokhai/scraper/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	redisHost := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		// MinIdleConns: 200,
		// PoolSize: 12000,
		// PoolTimeout: 240,
		Password: cfg.Redis.Password,
		DB: cfg.Redis.Db,	
	})

	fmt.Println("Connected to Redis successfully")
	return rdb
}