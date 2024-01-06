package usecase

import (
	"context"
	"strings"

	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/auth"
	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/pkg/utils"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

type authUseCase struct {
	userRepo       	auth.UserRepository
	cfg 			*config.Config
}

func NewAuthUseCase(
	cfg *config.Config,
	userRepo auth.UserRepository,
	) auth.UseCase {
	return &authUseCase{
		cfg:			cfg,
		userRepo:       userRepo,
	}
}

func (a *authUseCase) SignUp(ctx context.Context, username, password string) (*models.UserWithToken, error) {
	fmtusername := strings.ToLower(username)
	existingUser, _ := a.userRepo.GetUserByUsername(ctx, fmtusername)

	if existingUser != nil {
		return nil, auth.ErrUserExisted
	}

	user := &models.User{
		UserID: uuid.New().String(),
		Username: fmtusername,
		Password: password,
	}

	user.HashPassword()
	err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWTToken(user, a.cfg)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User: user,
		Token: token,
	}, nil
}

func (a *authUseCase) SignIn(ctx context.Context, username, password string) (*models.UserWithToken, error) {
	user, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, auth.ErrUserNotFound
	}

	if !user.ComparePassword(password) {
		return nil, auth.ErrWrongPassword
	}

	token, err := utils.GenerateJWTToken(user, a.cfg)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User: user,
		Token: token,
	}, nil
}

func (u *authUseCase) GetByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := u.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}