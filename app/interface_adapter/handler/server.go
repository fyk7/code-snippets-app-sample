package interface_adapter

import (
	"net/http"

	"github.com/fyk7/code-snippets-app/app/interface_adapter/handler/generated"
	"github.com/fyk7/code-snippets-app/app/usecase"
	"github.com/labstack/echo/v4"
)

// server implements generated.ServerInterface.
// This is the bridge between the OpenAPI-generated contract and the clean architecture usecase layer.
type server struct {
	snippetService usecase.SnippetService
	tagService     usecase.TagService
}

// NewServer creates a handler that satisfies the generated ServerInterface.
func NewServer(snippetService usecase.SnippetService, tagService usecase.TagService) generated.ServerInterface {
	return &server{
		snippetService: snippetService,
		tagService:     tagService,
	}
}

func (s *server) GetSnippetByID(c echo.Context, snippetId int) error {
	ctx := c.Request().Context()
	snippet, err := s.snippetService.GetByID(ctx, uint64(snippetId))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.ToSnippetResponse(snippet))
}

func (s *server) SearchSnippets(c echo.Context, params generated.SearchSnippetsParams) error {
	ctx := c.Request().Context()
	keyword := ""
	if params.SnippetKeyword != nil {
		keyword = *params.SnippetKeyword
	}
	snippets, err := s.snippetService.GetByKeyWord(ctx, keyword)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.ToSnippetResponses(snippets))
}

func (s *server) ListSnippetsByTag(c echo.Context, tagId int) error {
	ctx := c.Request().Context()
	snippets, err := s.snippetService.GetByKeyTagID(ctx, uint64(tagId))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.ToSnippetResponses(snippets))
}

func (s *server) CreateSnippet(c echo.Context) error {
	var req generated.SnippetCreateRequest
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0) // dummy user
	if err := s.snippetService.Create(ctx, generated.SnippetCreateToModel(req), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, generated.OKResponse{Messages: []string{"Successfully Created."}})
}

func (s *server) UpdateSnippet(c echo.Context) error {
	var req generated.SnippetUpdateRequest
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0)
	if err := s.snippetService.Update(ctx, generated.SnippetUpdateToModel(req), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.OKResponse{Messages: []string{"Successfully Updated."}})
}

func (s *server) AssociateSnippetWithTag(c echo.Context) error {
	var req generated.AssociateWithTagRequest
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := int64(0)
	if err := s.snippetService.AssociateWithTag(ctx, req.SnippetId, req.TagId, userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, generated.OKResponse{Messages: []string{"Successfully Created."}})
}

func (s *server) GetTagByID(c echo.Context, tagId int) error {
	ctx := c.Request().Context()
	tag, err := s.tagService.GetByID(ctx, uint64(tagId))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.ToTagResponse(tag))
}

func (s *server) SearchTags(c echo.Context, params generated.SearchTagsParams) error {
	ctx := c.Request().Context()
	keyword := ""
	if params.TagKeyword != nil {
		keyword = *params.TagKeyword
	}
	tags, err := s.tagService.GetByKeyWord(ctx, keyword)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, generated.ToTagResponses(tags))
}

func (s *server) CreateTag(c echo.Context) error {
	var req generated.TagCreateRequest
	if err := c.Bind(&req); err != nil {
		return handleError(c, err)
	}
	ctx := c.Request().Context()
	userID := uint64(0)
	if err := s.tagService.Create(ctx, generated.TagCreateToModel(req), userID); err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, generated.OKResponse{Messages: []string{"Successfully Created."}})
}
