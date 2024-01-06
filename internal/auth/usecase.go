package auth

import (
	"context"
	"github.com/pcaokhai/scraper/internal/models"
)

type UseCase interface {
	SignUp(ctx context.Context, username, password string) (*models.UserWithToken, error)
	SignIn(ctx context.Context, username, password string) (*models.UserWithToken, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
}
