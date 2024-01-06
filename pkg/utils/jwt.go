package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/models"
)

type Claims struct {
	Username 	string 		`json:"username"`
	ID    		string 		`json:"id"`
	jwt.RegisteredClaims
}


func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 60)
	claims := &Claims{
		Username: user.Username,
		ID:    user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

