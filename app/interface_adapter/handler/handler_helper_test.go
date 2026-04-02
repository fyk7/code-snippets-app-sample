package interface_adapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func callHandleError(err error) (int, ErrorResponseBody) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handleError(c, err)

	var body ErrorResponseBody
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	return rec.Code, body
}

func TestHandleError_NotFound(t *testing.T) {
	code, body := callHandleError(model.ErrNotFound)
	assert.Equal(t, http.StatusNotFound, code)
	assert.Contains(t, body.Messages[0], "not found")
}

func TestHandleError_Conflict(t *testing.T) {
	code, body := callHandleError(model.ErrConflict)
	assert.Equal(t, http.StatusConflict, code)
	assert.Contains(t, body.Messages[0], "already exists")
}

func TestHandleError_BadParam(t *testing.T) {
	code, _ := callHandleError(model.ErrBadParamInput)
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestHandleError_UnknownError(t *testing.T) {
	code, body := callHandleError(errors.New("unexpected"))
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Contains(t, body.Messages[0], "internal server error")
}

func TestHandleError_ValidationError(t *testing.T) {
	err := ValidRequest(struct {
		Name string `validate:"required"`
	}{})

	code, body := callHandleError(err)
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Contains(t, body.Messages[0], "validation error")
}
