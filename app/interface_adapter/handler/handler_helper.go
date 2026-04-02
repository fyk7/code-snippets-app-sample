package interface_adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OKResponseBody struct {
	Messages []string `json:"messages"`
}

type ErrorResponseBody struct {
	Messages []string `json:"messages"`
}

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return c.JSON(http.StatusNotFound, ErrorResponseBody{Messages: []string{model.ErrNotFound.Error()}})
	case errors.Is(err, model.ErrConflict):
		return c.JSON(http.StatusConflict, ErrorResponseBody{Messages: []string{model.ErrConflict.Error()}})
	case errors.Is(err, model.ErrBadParamInput):
		return c.JSON(http.StatusBadRequest, ErrorResponseBody{Messages: []string{model.ErrBadParamInput.Error()}})
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		type validationErr struct {
			Field        string `json:"field"`
			InvalidValue any    `json:"invalid_value"`
			Tag          string `json:"tag"`
			Param        string `json:"param"`
		}
		var validationErrs []validationErr
		for _, e := range ve {
			validationErrs = append(validationErrs, validationErr{
				Field:        e.Field(),
				InvalidValue: e.Value(),
				Tag:          e.Tag(),
				Param:        e.Param(),
			})
		}
		b, marshalErr := json.Marshal(validationErrs)
		if marshalErr != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponseBody{Messages: []string{model.ErrInternalServerError.Error()}})
		}
		return c.JSON(http.StatusBadRequest, ErrorResponseBody{Messages: []string{fmt.Sprintf("validation error: %s", string(b))}})
	}

	return c.JSON(http.StatusInternalServerError, ErrorResponseBody{Messages: []string{model.ErrInternalServerError.Error()}})
}
