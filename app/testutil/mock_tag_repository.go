package testutil

import (
	"context"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) GetAll(ctx context.Context) ([]model.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagRepository) GetByID(ctx context.Context, tagID uint64) (model.Tag, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).(model.Tag), args.Error(1)
}

func (m *MockTagRepository) FindByKeyWord(ctx context.Context, keyword string) ([]model.Tag, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagRepository) Create(ctx context.Context, tag model.Tag, userID uint64) error {
	args := m.Called(ctx, tag, userID)
	return args.Error(0)
}
