package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestSnippetService(repo *testutil.MockSnippetRepository) SnippetService {
	return NewSnippetService(repo, 5*time.Second)
}

func TestSnippetService_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		expected := []model.Snippet{
			{SnippetID: 1, Title: "Test", Body: "body", ProgramingLanguage: "go"},
			{SnippetID: 2, Title: "Test2", Body: "body2", ProgramingLanguage: "python"},
		}
		repo.On("GetAll", mock.Anything).Return(expected, nil)

		result, err := svc.List(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		repo.On("GetAll", mock.Anything).Return([]model.Snippet(nil), errors.New("db error"))

		result, err := svc.List(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestSnippetService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		expected := model.Snippet{SnippetID: 1, Title: "Test", Body: "body", ProgramingLanguage: "go"}
		repo.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		result, err := svc.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		repo.On("GetByID", mock.Anything, uint64(999)).Return(model.Snippet{}, model.ErrNotFound)

		_, err := svc.GetByID(context.Background(), 999)

		assert.ErrorIs(t, err, model.ErrNotFound)
		repo.AssertExpectations(t)
	})
}

func TestSnippetService_GetByKeyWord(t *testing.T) {
	repo := new(testutil.MockSnippetRepository)
	svc := newTestSnippetService(repo)

	expected := []model.Snippet{
		{SnippetID: 1, Title: "Go Error Handling", Body: "body"},
	}
	repo.On("FindByKeyWord", mock.Anything, "error").Return(expected, nil)

	result, err := svc.GetByKeyWord(context.Background(), "error")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestSnippetService_GetByKeyTagID(t *testing.T) {
	repo := new(testutil.MockSnippetRepository)
	svc := newTestSnippetService(repo)

	expected := []model.Snippet{
		{SnippetID: 1, Title: "Tagged Snippet"},
	}
	repo.On("FindByTag", mock.Anything, uint64(5)).Return(expected, nil)

	result, err := svc.GetByKeyTagID(context.Background(), 5)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestSnippetService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		snippet := model.Snippet{Title: "New", Body: "body", ProgramingLanguage: "go"}
		repo.On("Create", mock.Anything, snippet, uint64(1)).Return(nil)

		err := svc.Create(context.Background(), snippet, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(testutil.MockSnippetRepository)
		svc := newTestSnippetService(repo)

		snippet := model.Snippet{Title: "New", Body: "body"}
		repo.On("Create", mock.Anything, snippet, uint64(1)).Return(errors.New("db error"))

		err := svc.Create(context.Background(), snippet, 1)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestSnippetService_Update(t *testing.T) {
	repo := new(testutil.MockSnippetRepository)
	svc := newTestSnippetService(repo)

	snippet := model.Snippet{SnippetID: 1, Title: "Updated", Body: "body"}
	repo.On("Update", mock.Anything, snippet, uint64(1)).Return(nil)

	err := svc.Update(context.Background(), snippet, 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSnippetService_Delete(t *testing.T) {
	repo := new(testutil.MockSnippetRepository)
	svc := newTestSnippetService(repo)

	repo.On("Delete", mock.Anything, uint64(1)).Return(nil)

	err := svc.Delete(context.Background(), 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSnippetService_AssociateWithTag(t *testing.T) {
	repo := new(testutil.MockSnippetRepository)
	svc := newTestSnippetService(repo)

	repo.On("AssociateWithTag", mock.Anything, int64(1), int64(5), int64(0)).Return(nil)

	err := svc.AssociateWithTag(context.Background(), 1, 5, 0)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
