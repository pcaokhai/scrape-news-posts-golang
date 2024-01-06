package delivery

import (
	"github.com/labstack/echo/v4"

	"github.com/pcaokhai/scraper/internal/auth"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handler) {
	authGroup.POST("/signup", h.SignUp())
	authGroup.POST("/signin", h.SignIn())
	authGroup.POST("/logout", h.Logout())
}