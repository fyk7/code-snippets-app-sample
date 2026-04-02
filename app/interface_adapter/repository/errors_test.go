package interface_adapter

import (
	"errors"
	"testing"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToDomainError(t *testing.T) {
	t.Run("nil error returns nil", func(t *testing.T) {
		assert.NoError(t, toDomainError(nil))
	})

	t.Run("gorm ErrRecordNotFound maps to ErrNotFound", func(t *testing.T) {
		err := toDomainError(gorm.ErrRecordNotFound)
		assert.ErrorIs(t, err, model.ErrNotFound)
	})

	t.Run("duplicate entry maps to ErrConflict", func(t *testing.T) {
		err := toDomainError(errors.New("Error 1062: Duplicate entry 'test' for key 'PRIMARY'"))
		assert.ErrorIs(t, err, model.ErrConflict)
	})

	t.Run("unknown error passes through", func(t *testing.T) {
		original := errors.New("some unknown error")
		err := toDomainError(original)
		assert.Equal(t, original, err)
	})
}
