package testutil

import (
	"context"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockSnippetService struct {
	mock.Mock
}

func (m *MockSnippetService) List(ctx context.Context) ([]model.Snippet, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetService) GetByID(ctx context.Context, id uint64) (model.Snippet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Snippet), args.Error(1)
}

func (m *MockSnippetService) GetByKeyWord(ctx context.Context, keyword string) ([]model.Snippet, error) {
	args := m.Called(ctx, keyword)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetService) GetByKeyTagID(ctx context.Context, tagID uint64) ([]model.Snippet, error) {
	args := m.Called(ctx, tagID)
	return args.Get(0).([]model.Snippet), args.Error(1)
}

func (m *MockSnippetService) AssociateWithTag(ctx context.Context, snippetID, tagID, userID int64) error {
	args := m.Called(ctx, snippetID, tagID, userID)
	return args.Error(0)
}

func (m *MockSnippetService) Create(ctx context.Context, snippet model.Snippet, userID uint64) error {
	args := m.Called(ctx, snippet, userID)
	return args.Error(0)
}

func (m *MockSnippetService) Update(ctx context.Context, snippet model.Snippet, userID uint64) error {
	args := m.Called(ctx, snippet, userID)
	return args.Error(0)
}

func (m *MockSnippetService) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
