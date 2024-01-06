package server

import (
	"github.com/labstack/echo/v4"
	postHandler "github.com/pcaokhai/scraper/internal/post/handler"
	postRepository "github.com/pcaokhai/scraper/internal/post/repository"
	postUsecase "github.com/pcaokhai/scraper/internal/post/usecase"

	authHandler "github.com/pcaokhai/scraper/internal/auth/delivery"
	authRepository "github.com/pcaokhai/scraper/internal/auth/repository"
	authUsecase "github.com/pcaokhai/scraper/internal/auth/usecase"

	templateHandler "github.com/pcaokhai/scraper/internal/template/handlers"

	middleware "github.com/pcaokhai/scraper/internal/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// repo
	postRepo := postRepository.NewPostRepository(s.db)
	postRedisRepo := postRepository.NewPostsRedisRepo(s.redis)
	authRepo := authRepository.NewDBRepository(s.db)

	//usecase
	postUC := postUsecase.NewPostUseCase(postRepo, postRedisRepo)
	authUC := authUsecase.NewAuthUseCase(s.cfg, authRepo)

	//middleware
	mw := middleware.NewMiddlewareManager(s.cfg, authUC)

	//handler 
	postHler := postHandler.NewPostHandler(postUC)
	authHler := authHandler.NewAuthHandler(s.cfg, authUC)

	v1 := e.Group("/api/v1")
	postGroup := v1.Group("/post")
	authGroup := v1.Group("/auth")

	tmplHandler := templateHandler.NewTemplateHandler(s.db, s.redis, postRedisRepo)

	e.GET("/login", tmplHandler.SignInPage()).Name="loginForm"
	e.GET("/register", tmplHandler.SignUpPage()).Name="registerForm"
	e.GET("/admin", tmplHandler.AdminPage(), mw.AuthJWTMiddleware()).Name="adminPage"
	e.GET("/", tmplHandler.HomePage()).Name="homePage"

	postHandler.MapPostRoutes(postGroup, postHler)
	authHandler.MapAuthRoutes(authGroup, authHler)

	return nil
}