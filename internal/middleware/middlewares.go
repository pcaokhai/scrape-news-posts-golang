package middleware

import (
	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/auth"
)

type MiddlewareManager struct {
	cfg     *config.Config
	authUC  auth.UseCase
}

func NewMiddlewareManager(cfg *config.Config, authUC auth.UseCase) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, authUC: authUC}
}