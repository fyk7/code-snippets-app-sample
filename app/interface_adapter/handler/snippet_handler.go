package interface_adapter

import (
	"net/http"
	"strconv"

	"github.com/fyk7/code-snippets-app/app/usecase"
	"github.com/labstack/echo/v4"
)

type snippetHandler struct {
	snippetService usecase.SnippetService
}

func NewSnippetHandler(e *echo.Echo, s usecase.SnippetService) {
	handler := &snippetHandler{
		snippetService: s,
	}
	e.GET("/snippets/:snippet_id", handler.GetSnippetByID)
	e.GET("/snippets/search", handler.FindSnippetByKeyWord)
	e.GET("/snippets/tags/:tag_id", handler.ListByTag)
	e.POST("/snippets", handler.PostSnippet)
	e.POST("/snippets/associate", handler.AssociateWithTag)
}

func (h *snippetHandler) GetSnippetByID(c echo.Context) error {
	ctx := c.Request().Context()
	snippetID, err := strconv.Atoi(c.Param("snippet_id"))
	if err != nil {
		return handleError(c, err)
	}
	snippet, err := h.snippetService.GetByID(ctx, uint64(snippetID))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, snippet)
}

func (h *snippetHandler) FindSnippetByKeyWord(c echo.Context) error {
	ctx := c.Request().Context()
	snippetKeyword := c.QueryParam("snippet_keyword")
	snippet, err := h.snippetService.GetByKeyWord(ctx, snippetKeyword)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, snippet)
}

func (h *snippetHandler) ListByTag(c echo.Context) error {
	ctx := c.Request().Context()
	tagID, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		return handleError(c, err)
	}
	snippets, err := h.snippetService.GetByKeyTagID(ctx, uint64(tagID))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, snippets)
}

func (h *snippetHandler) PostSnippet(c echo.Context) error {
	var req SnippetPostReq
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	if err := ValidRequest(req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0) // dummy user
	if err := h.snippetService.Create(ctx, req.ConvertToModel(), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, OKResponseBody{Messages: []string{"Successfully Created."}})
}

type AssociateWithTagReq struct {
	SnippetID int64 `json:"snippet_id" validate:"required"`
	TagID     int64 `json:"tag_id" validate:"required"`
}

func (h *snippetHandler) AssociateWithTag(c echo.Context) error {
	var req AssociateWithTagReq
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	if err := ValidRequest(req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := int64(0)
	if err := h.snippetService.AssociateWithTag(ctx, req.SnippetID, req.TagID, userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, OKResponseBody{Messages: []string{"Successfully Created."}})
}

func (h *snippetHandler) PutSnippet(c echo.Context) error {
	var req SnippetPutReq
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	if err := ValidRequest(req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0)
	if err := h.snippetService.Update(ctx, req.ConvertToModel(), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, OKResponseBody{Messages: []string{"Successfully Updated."}})
}
