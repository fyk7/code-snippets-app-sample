package interface_adapter

import (
	"net/http"
	"strconv"

	"github.com/fyk7/code-snippets-app/app/usecase"
	"github.com/labstack/echo/v4"
)

type tagHandler struct {
	tagService usecase.TagService
}

func NewTagHandler(e *echo.Echo, s usecase.TagService) {
	handler := &tagHandler{
		tagService: s,
	}
	e.GET("/tags/:tag_id", handler.GetTagByID)
	e.GET("/tags/search", handler.FindTagByKeyWord)
	e.POST("/tags", handler.PostTag)
}

func (h *tagHandler) GetTagByID(c echo.Context) error {
	ctx := c.Request().Context()
	tagID, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		return handleError(c, err)
	}
	tag, err := h.tagService.GetByID(ctx, uint64(tagID))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, tag)
}

func (h *tagHandler) FindTagByKeyWord(c echo.Context) error {
	ctx := c.Request().Context()
	tagKeyword := c.QueryParam("tag_keyword")
	tags, err := h.tagService.GetByKeyWord(ctx, tagKeyword)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, tags)
}

func (h *tagHandler) PostTag(c echo.Context) error {
	var req TagPostReq
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	if err := ValidRequest(req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0)
	if err := h.tagService.Create(ctx, req.ConvertToModel(), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, OKResponseBody{Messages: []string{"Successfully Created."}})
}
