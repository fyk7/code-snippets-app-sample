package interface_adapter

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidRequest(req any) error {
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}
	return nil
}
