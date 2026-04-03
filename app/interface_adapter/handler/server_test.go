package interface_adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/interface_adapter/handler/generated"
	"github.com/fyk7/code-snippets-app/app/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupServer(snippetSvc *testutil.MockSnippetService, tagSvc *testutil.MockTagService) *echo.Echo {
	e := echo.New()
	srv := NewServer(snippetSvc, tagSvc)
	generated.RegisterHandlers(e, srv)
	return e
}

// --- Snippet endpoints ---

func TestServer_GetSnippetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		expected := model.Snippet{SnippetID: 1, Title: "Test", Body: "body", ProgramingLanguage: "go"}
		snippetSvc.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/snippets/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result generated.Snippet
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result.Title)
		assert.Equal(t, uint64(1), result.SnipetId)
		snippetSvc.AssertExpectations(t)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		req := httptest.NewRequest(http.MethodGet, "/snippets/abc", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		snippetSvc.On("GetByID", mock.Anything, uint64(999)).Return(model.Snippet{}, model.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/snippets/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		snippetSvc.AssertExpectations(t)
	})
}

func TestServer_SearchSnippets(t *testing.T) {
	snippetSvc := new(testutil.MockSnippetService)
	tagSvc := new(testutil.MockTagService)
	e := setupServer(snippetSvc, tagSvc)

	expected := []model.Snippet{
		{SnippetID: 1, Title: "Go Error Handling"},
	}
	snippetSvc.On("GetByKeyWord", mock.Anything, "error").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/snippets/search?snippet_keyword=error", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []generated.Snippet
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	snippetSvc.AssertExpectations(t)
}

func TestServer_ListSnippetsByTag(t *testing.T) {
	snippetSvc := new(testutil.MockSnippetService)
	tagSvc := new(testutil.MockTagService)
	e := setupServer(snippetSvc, tagSvc)

	expected := []model.Snippet{{SnippetID: 1, Title: "Tagged"}}
	snippetSvc.On("GetByKeyTagID", mock.Anything, uint64(5)).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/snippets/tags/5", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	snippetSvc.AssertExpectations(t)
}

func TestServer_CreateSnippet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		snippetSvc.On("Create", mock.Anything, mock.AnythingOfType("model.Snippet"), uint64(0)).Return(nil)

		body := `{"title":"Test","body":"code","programing_language":"go"}`
		req := httptest.NewRequest(http.MethodPost, "/snippets", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		snippetSvc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		req := httptest.NewRequest(http.MethodPost, "/snippets", strings.NewReader("{invalid"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestServer_AssociateSnippetWithTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		snippetSvc.On("AssociateWithTag", mock.Anything, int64(1), int64(5), int64(0)).Return(nil)

		body := `{"snippet_id":1,"tag_id":5}`
		req := httptest.NewRequest(http.MethodPost, "/snippets/associate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		snippetSvc.AssertExpectations(t)
	})

	t.Run("conflict", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		snippetSvc.On("AssociateWithTag", mock.Anything, int64(1), int64(5), int64(0)).Return(model.ErrConflict)

		body := `{"snippet_id":1,"tag_id":5}`
		req := httptest.NewRequest(http.MethodPost, "/snippets/associate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		snippetSvc.AssertExpectations(t)
	})
}

// --- Tag endpoints ---

func TestServer_GetTagByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		expected := model.Tag{TagID: 1, TagName: "golang"}
		tagSvc.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result generated.Tag
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "golang", result.TagName)
		tagSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		tagSvc.On("GetByID", mock.Anything, uint64(999)).Return(model.Tag{}, model.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/tags/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		tagSvc.AssertExpectations(t)
	})
}

func TestServer_SearchTags(t *testing.T) {
	snippetSvc := new(testutil.MockSnippetService)
	tagSvc := new(testutil.MockTagService)
	e := setupServer(snippetSvc, tagSvc)

	expected := []model.Tag{{TagID: 1, TagName: "golang"}}
	tagSvc.On("GetByKeyWord", mock.Anything, "go").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/tags/search?tag_keyword=go", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []generated.Tag
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	tagSvc.AssertExpectations(t)
}

func TestServer_CreateTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		snippetSvc := new(testutil.MockSnippetService)
		tagSvc := new(testutil.MockTagService)
		e := setupServer(snippetSvc, tagSvc)

		tagSvc.On("Create", mock.Anything, mock.AnythingOfType("model.Tag"), uint64(0)).Return(nil)

		body := `{"tag_name":"rust"}`
		req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		tagSvc.AssertExpectations(t)
	})
}
