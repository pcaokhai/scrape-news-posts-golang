package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pcaokhai/scraper/internal/post"
)

func MapPostRoutes(postGroup *echo.Group, handler post.PostHandler) {
	postGroup.GET("/", handler.GetAllPosts())
	postGroup.PUT("/:postId", handler.UpdatePost())
	postGroup.DELETE("/:postId", handler.DeletePost())
}