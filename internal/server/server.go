package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pcaokhai/scraper/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	cfg 	*config.Config
	echo 	*echo.Echo
	db		*gorm.DB
	ready 	chan bool
	redis 	*redis.Client
}

func NewServer(
		cfg *config.Config,
		db *gorm.DB, 
		redis *redis.Client, 
		ready chan bool, 
		) *Server {
	return &Server{
		cfg: cfg,
		echo: echo.New(), 
		db: db, 
		redis: redis,
		ready: ready,
	}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr: ":" + s.cfg.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	go func() {
		fmt.Printf("Server is listening on PORT: %s", s.cfg.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Error starting Server: %v", err)
		}
	}()

	// set up routes
	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	// gracefully shut down server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down")
	ctx, shutdown := context.WithTimeout(context.Background(), 15 * time.Second)
	defer func() {
		dbInstance, _ := s.db.DB()
		dbInstance.Close()

		s.redis.Close()
		shutdown()
	}()
	log.Println("Stopping http server")

	if err := s.echo.Server.Shutdown(ctx); err == context.DeadlineExceeded {
		log.Println("Halted active connection")
		return err
	}

	close(quit)
	return nil
}