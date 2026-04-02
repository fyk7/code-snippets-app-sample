package di

import (
	"time"

	"github.com/fyk7/code-snippets-app/app/config"
	"github.com/fyk7/code-snippets-app/app/infrastructure/database"
	repository "github.com/fyk7/code-snippets-app/app/interface_adapter/repository"
	"github.com/fyk7/code-snippets-app/app/usecase"
	"gorm.io/gorm"
)

type ServiceContainer struct {
	SnippetService usecase.SnippetService
	TagService     usecase.TagService
	UserService    usecase.UserService
	DB             *gorm.DB
}

func Initialize(cfg *config.Config, timeout time.Duration) *ServiceContainer {
	db := database.NewDB(cfg)

	snippetRepo := repository.NewSnippetRepository(db)
	tagRepo := repository.NewTagRepository(db)
	userRepo := repository.NewUserRepository(db)

	snippetSvc := usecase.NewSnippetService(snippetRepo, timeout)
	tagSvc := usecase.NewTagService(tagRepo, timeout)
	userSvc := usecase.NewUserService(userRepo, timeout)

	return &ServiceContainer{
		SnippetService: snippetSvc,
		TagService:     tagSvc,
		UserService:    userSvc,
		DB:             db,
	}
}
