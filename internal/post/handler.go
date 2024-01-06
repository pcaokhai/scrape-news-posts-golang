package post

import (
	"github.com/labstack/echo/v4"
)

type PostHandler interface {
	GetAllPosts() echo.HandlerFunc
	UpdatePost() echo.HandlerFunc
	DeletePost() echo.HandlerFunc
}