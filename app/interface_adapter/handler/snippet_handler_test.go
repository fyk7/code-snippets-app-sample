package interface_adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSnippetHandler(mockSvc *testutil.MockSnippetService) *echo.Echo {
	e := echo.New()
	NewSnippetHandler(e, mockSvc)
	return e
}

func TestGetSnippetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		expected := model.Snippet{SnippetID: 1, Title: "Test", Body: "body", ProgramingLanguage: "go"}
		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/snippets/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result model.Snippet
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "Test", result.Title)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/snippets/abc", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		mockSvc.On("GetByID", mock.Anything, uint64(999)).Return(model.Snippet{}, model.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/snippets/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestFindSnippetByKeyWord(t *testing.T) {
	mockSvc := new(testutil.MockSnippetService)
	e := setupSnippetHandler(mockSvc)

	expected := []model.Snippet{
		{SnippetID: 1, Title: "Go Error Handling"},
	}
	mockSvc.On("GetByKeyWord", mock.Anything, "error").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/snippets/search?snippet_keyword=error", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []model.Snippet
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockSvc.AssertExpectations(t)
}

func TestListByTag(t *testing.T) {
	mockSvc := new(testutil.MockSnippetService)
	e := setupSnippetHandler(mockSvc)

	expected := []model.Snippet{
		{SnippetID: 1, Title: "Tagged"},
	}
	mockSvc.On("GetByKeyTagID", mock.Anything, uint64(5)).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/snippets/tags/5", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockSvc.AssertExpectations(t)
}

func TestPostSnippet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		mockSvc.On("Create", mock.Anything, mock.AnythingOfType("model.Snippet"), uint64(0)).Return(nil)

		body := `{"title":"Test","body":"code","programing_language":"go"}`
		req := httptest.NewRequest(http.MethodPost, "/snippets", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("validation error - missing required fields", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		body := `{"description":"only description"}`
		req := httptest.NewRequest(http.MethodPost, "/snippets", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/snippets", strings.NewReader("{invalid"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAssociateWithTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		mockSvc.On("AssociateWithTag", mock.Anything, int64(1), int64(5), int64(0)).Return(nil)

		body := `{"snippet_id":1,"tag_id":5}`
		req := httptest.NewRequest(http.MethodPost, "/snippets/associate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("conflict error", func(t *testing.T) {
		mockSvc := new(testutil.MockSnippetService)
		e := setupSnippetHandler(mockSvc)

		mockSvc.On("AssociateWithTag", mock.Anything, int64(1), int64(5), int64(0)).Return(model.ErrConflict)

		body := `{"snippet_id":1,"tag_id":5}`
		req := httptest.NewRequest(http.MethodPost, "/snippets/associate", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}
