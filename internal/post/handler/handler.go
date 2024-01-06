package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pcaokhai/scraper/internal/models"
	"github.com/pcaokhai/scraper/internal/post"
	"github.com/pcaokhai/scraper/internal/post/presenter"
)
const CtxPostId = "postId"

type postHandler struct {
	usecase post.PostUseCase
}

func NewPostHandler(usecase post.PostUseCase) post.PostHandler {
	return &postHandler{usecase: usecase,}
}

func (handler *postHandler) GetAllPosts() echo.HandlerFunc{
	return func(c echo.Context) error {
		posts, err := handler.usecase.GetAllPosts(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		response := map[string]interface{} {
			"msg": "Fetched all posts successfully",
			"data": mapPosts(posts),
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (handler *postHandler) UpdatePost() echo.HandlerFunc{
	return func(c echo.Context) error {
		postId := c.Param("postId")
		id, _ := strconv.Atoi(postId)

		postUpdateRequest := new(presenter.PostUpdateRequest)
		if err := c.Bind(postUpdateRequest); err != nil {
			data := map[string]interface{}{
            "message": err.Error(),
        }
        	return c.JSON(http.StatusBadRequest, data)
		}

		updatedPost, err := handler.usecase.UpdatePost(c.Request().Context(), id, postUpdateRequest.Title)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		response := map[string]interface{} {
			"msg": "Post updated successfully",
			"data": updatedPost,
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (handler *postHandler) DeletePost() echo.HandlerFunc{
	return func(c echo.Context) error {
		postId := c.Param("postId")
		id, _ := strconv.Atoi(postId)

		err := handler.usecase.DeletePost(c.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		response := map[string]interface{} {
			"msg": "Post deleted successfully",
		}

		return c.JSON(http.StatusOK, response)
	}
}

func mapPosts(posts []*models.Post) []*presenter.PostResponse {
	output := make([]*presenter.PostResponse, len(posts))
	for i, v := range posts {
		output[i] = mapPost(v)
	}

	return output
}

func mapPost(post *models.Post) *presenter.PostResponse {
	return &presenter.PostResponse {
		Id: post.Id,
		Title: post.Title,
		PostId: post.PostId,
		Url: post.Url,
	}
}