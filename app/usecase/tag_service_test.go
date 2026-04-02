package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestTagService(repo *testutil.MockTagRepository) TagService {
	return NewTagService(repo, 5*time.Second)
}

func TestTagService_List(t *testing.T) {
	repo := new(testutil.MockTagRepository)
	svc := newTestTagService(repo)

	expected := []model.Tag{
		{TagID: 1, TagName: "golang"},
		{TagID: 2, TagName: "python"},
	}
	repo.On("GetAll", mock.Anything).Return(expected, nil)

	result, err := svc.List(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestTagService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(testutil.MockTagRepository)
		svc := newTestTagService(repo)

		expected := model.Tag{TagID: 1, TagName: "golang"}
		repo.On("GetByID", mock.Anything, uint64(1)).Return(expected, nil)

		result, err := svc.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, "golang", result.TagName)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(testutil.MockTagRepository)
		svc := newTestTagService(repo)

		repo.On("GetByID", mock.Anything, uint64(999)).Return(model.Tag{}, model.ErrNotFound)

		_, err := svc.GetByID(context.Background(), 999)

		assert.ErrorIs(t, err, model.ErrNotFound)
		repo.AssertExpectations(t)
	})
}

func TestTagService_GetByKeyWord(t *testing.T) {
	repo := new(testutil.MockTagRepository)
	svc := newTestTagService(repo)

	expected := []model.Tag{{TagID: 1, TagName: "golang"}}
	repo.On("FindByKeyWord", mock.Anything, "go").Return(expected, nil)

	result, err := svc.GetByKeyWord(context.Background(), "go")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

func TestTagService_Create(t *testing.T) {
	repo := new(testutil.MockTagRepository)
	svc := newTestTagService(repo)

	tag := model.Tag{TagName: "rust"}
	repo.On("Create", mock.Anything, tag, uint64(1)).Return(nil)

	err := svc.Create(context.Background(), tag, 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
