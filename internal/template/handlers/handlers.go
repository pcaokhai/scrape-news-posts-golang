package handlers

import (
	"html/template"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/pcaokhai/scraper/pkg/utils"

	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/template/presenter"
	"github.com/pcaokhai/scraper/internal/post"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type templateHandler struct {
	db *gorm.DB
	rdb *redis.Client
	redisRepo post.PostRedisRepository
}

func NewTemplateHandler(db *gorm.DB, rdb *redis.Client, redisRepo post.PostRedisRepository) *templateHandler {
	return &templateHandler{db: db, rdb: rdb, redisRepo: redisRepo}
}

func (t *templateHandler) SignInPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		fp := path.Join("templates", "signin.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (t *templateHandler) SignUpPage() echo.HandlerFunc {
	return func(c echo.Context) error {

		fp := path.Join("templates", "signup.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (t *templateHandler) HomePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		posts, err := utils.FetchData(c.Request().Context(), t.db, t.rdb, t.redisRepo)
		postsData := mapPosts(posts)
		
		fp := path.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := tmpl.Execute(c.Response().Writer, postsData); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func (t *templateHandler) AdminPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		fp := path.Join("templates", "admin.html")
		tmpl, err := template.ParseFiles(fp)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

func mapPosts(posts []*models.Post) []*presenter.PostData {
	output := make([]*presenter.PostData, len(posts))
	for i, v := range posts {
		output[i] = mapPost(v)
	}

	return output
}

func mapPost(post *models.Post) *presenter.PostData {
	return &presenter.PostData {
		Title: post.Title,
		Url: post.Url,
		ArticleId: post.PostId,
	}
}