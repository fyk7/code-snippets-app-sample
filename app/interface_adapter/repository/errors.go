package interface_adapter

import (
	"errors"
	"strings"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"gorm.io/gorm"
)

// toDomainError converts GORM errors to domain errors.
func toDomainError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrNotFound
	}
	// MySQL duplicate entry error code 1062
	if strings.Contains(err.Error(), "Duplicate entry") {
		return model.ErrConflict
	}
	return err
}
