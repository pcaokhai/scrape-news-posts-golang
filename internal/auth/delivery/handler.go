package delivery

import (
	"net/http"

	"github.com/pcaokhai/scraper/config"
	"github.com/pcaokhai/scraper/internal/auth"
	"github.com/pcaokhai/scraper/internal/auth/presenter"
	"github.com/pcaokhai/scraper/pkg/utils"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	cfg 			*config.Config
	authUC 			auth.UseCase
}

func NewAuthHandler( 
		cfg *config.Config,
		useCase auth.UseCase, 
	) auth.Handler {
	return &authHandler{
		cfg: cfg,
		authUC: useCase,
	}
}

func (h *authHandler) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &presenter.SignUpInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user, err := h.authUC.SignUp(c.Request().Context(), input.Username, input.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.SetCookie(utils.CreateJWTCookie(h.cfg, user.Token))

		return c.Redirect(http.StatusMovedPermanently, "/admin")
	}
}

func (h *authHandler) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &presenter.LoginInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		userWithToken, err := h.authUC.SignIn(c.Request().Context(), input.Username, input.Password)
		if err != nil {
			if err == auth.ErrUserNotFound {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			if err == auth.ErrWrongPassword {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		c.SetCookie(utils.CreateJWTCookie(h.cfg, userWithToken.Token))

		return c.Redirect(http.StatusMovedPermanently, "/admin")
	}
}

func (h *authHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie(h.cfg.CookieName)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "No cookie",
			})
		}

		utils.DeleteCookie(c, h.cfg.CookieName)

		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}