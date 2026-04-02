package testutil

import (
	"context"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockTagService struct {
	mock.Mock
}

func (m *MockTagService) List(ctx context.Context) ([]model.Tag, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagService) GetByID(ctx context.Context, id uint64) (model.Tag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Tag), args.Error(1)
}

func (m *MockTagService) GetByKeyWord(ctx context.Context, keyword string) ([]model.Tag, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]model.Tag), args.Error(1)
}

func (m *MockTagService) Create(ctx context.Context, tag model.Tag, userID uint64) error {
	args := m.Called(ctx, tag, userID)
	return args.Error(0)
}
