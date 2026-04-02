package usecase

import (
	"context"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/domain/repository"
)

type UserService interface {
	List(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, userID uint64) (model.User, error)
	GetByKeyWord(ctx context.Context, userName string) ([]model.User, error)
	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
}

type userService struct {
	repo           repository.UserRepository
	contextTimeout time.Duration
}

func NewUserService(repo repository.UserRepository, timeout time.Duration) UserService {
	return &userService{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (us *userService) List(ctx context.Context) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) GetByKeyWord(ctx context.Context, userName string) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	users, err := us.repo.FindByName(ctx, userName)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *userService) GetByID(ctx context.Context, userID uint64) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	user, err := us.repo.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (us *userService) Create(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	if err := us.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (us *userService) Update(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.contextTimeout)
	defer cancel()

	if err := us.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
