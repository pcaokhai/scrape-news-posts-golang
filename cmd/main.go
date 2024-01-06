package main

import (
	"fmt"
	"log"

	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/server"
	"github.com/pcaokhai/scraper/pkg/db/postgres"
	redis_conn "github.com/pcaokhai/scraper/pkg/db/redis"
	"github.com/pcaokhai/scraper/pkg/scraper"
)

func main () {
	// Load & parse configs
	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	fmt.Println("Configuarations loaded.")

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// connect postgres
	db, _ := postgres.NewPsqlDBConnection(cfg)
	err = db.AutoMigrate(&models.Post{}, &models.User{})
	if err != nil {
        panic(err)
    }

	// connect redis
	redisClient := redis_conn.NewRedisClient(cfg)

	scraper.Scrape(cfg)
	go scraper.ScheduledScrape(cfg)

	s := server.NewServer(cfg, db, redisClient, nil)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	
}