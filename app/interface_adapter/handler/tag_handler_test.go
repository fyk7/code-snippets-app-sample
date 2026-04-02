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

func setupTagHandler(mockSvc *testutil.MockTagService) *echo.Echo {
	e := echo.New()
	NewTagHandler(e, mockSvc)
	return e
}

func TestGetTagByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(testutil.MockTagService)
		e := setupTagHandler(mockSvc)

		expected := model.Tag{TagID: 1, TagName: "golang"}
		mockSvc.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/tags/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result model.Tag
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "golang", result.TagName)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(testutil.MockTagService)
		e := setupTagHandler(mockSvc)

		mockSvc.On("GetByID", mock.Anything, uint64(999)).Return(model.Tag{}, model.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/tags/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestFindTagByKeyWord(t *testing.T) {
	mockSvc := new(testutil.MockTagService)
	e := setupTagHandler(mockSvc)

	expected := []model.Tag{{TagID: 1, TagName: "golang"}}
	mockSvc.On("GetByKeyWord", mock.Anything, "go").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/tags/search?tag_keyword=go", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result []model.Tag
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockSvc.AssertExpectations(t)
}

func TestPostTag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(testutil.MockTagService)
		e := setupTagHandler(mockSvc)

		mockSvc.On("Create", mock.Anything, mock.AnythingOfType("model.Tag"), uint64(0)).Return(nil)

		body := `{"tag_name":"rust"}`
		req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("validation error - empty tag name", func(t *testing.T) {
		mockSvc := new(testutil.MockTagService)
		e := setupTagHandler(mockSvc)

		body := `{}`
		req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
