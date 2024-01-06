package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/auth"
	"github.com/pcaokhai/scraper/pkg/utils"
)

func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, mw.authUC, c, mw.cfg); err != nil {
					return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
				}

				return next(c)
			}

			cookie, err := c.Cookie(mw.cfg.CookieName)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, errors.New("c.Cookie"))
			}

			if err = mw.validateJWTToken(cookie.Value, mw.authUC, c, mw.cfg); err != nil {
				return c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
			}
			return next(c)
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return errors.New("Invalid token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return errors.New("Invalid token")
		}

		u, err := authUC.GetByID(c.Request().Context(), userID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}