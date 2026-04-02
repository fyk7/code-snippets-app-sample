package testutil

import (
	"context"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockSnippetRepository struct {
	mock.Mock
}

func (m *MockSnippetRepository) GetAll(ctx context.Context) ([]model.Snippet, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) GetByID(ctx context.Context, id uint64) (model.Snippet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) FindByKeyWord(ctx context.Context, keyword string) ([]model.Snippet, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) FindByTag(ctx context.Context, tagID uint64) ([]model.Snippet, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetRepository) AssociateWithTag(ctx context.Context, snippetID, tagID, userID int64) error {
	args := m.Called(ctx, snippetID, tagID, userID)
	return args.Error(0)
}

func (m *MockSnippetRepository) Create(ctx context.Context, s model.Snippet, userID uint64) error {
	args := m.Called(ctx, s, userID)
	return args.Error(0)
}

func (m *MockSnippetRepository) Update(ctx context.Context, s model.Snippet, userID uint64) error {
	args := m.Called(ctx, s, userID)
	return args.Error(0)
}

func (m *MockSnippetRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
