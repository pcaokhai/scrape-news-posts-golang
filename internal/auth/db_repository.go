package auth

import (
	"context"
	"github.com/pcaokhai/scraper/internal/models"
)

const CtxUserKey = "userId"

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserById(ctx context.Context, userId string) (*models.User, error)
}
