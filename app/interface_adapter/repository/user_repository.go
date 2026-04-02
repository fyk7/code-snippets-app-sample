package interface_adapter

import (
	"context"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		Conn: db,
	}
}

func (ur *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM user`
	if err := ur.Conn.WithContext(ctx).Raw(query).Scan(&users).Error; err != nil {
		return nil, toDomainError(err)
	}
	return users, nil
}

func (ur *userRepository) GetByID(ctx context.Context, userID uint64) (model.User, error) {
	var user model.User
	query := `
	SELECT *
	FROM user
	WHERE user_id = @userID
	`
	bindParams := map[string]any{
		"userID": userID,
	}
	if err := ur.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&user).Error; err != nil {
		return model.User{}, toDomainError(err)
	}
	return user, nil
}

func (ur *userRepository) FindByName(ctx context.Context, userName string) ([]model.User, error) {
	var users []model.User
	query := `
	SELECT *
	FROM user
	WHERE user_name LIKE CONCAT('%', @userName, '%')
	`
	bindParams := map[string]any{
		"userName": userName,
	}
	if err := ur.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&users).Error; err != nil {
		return nil, toDomainError(err)
	}
	return users, nil
}

func (ur *userRepository) Create(ctx context.Context, user model.User) error {
	query := `
	INSERT INTO user (
	  user_name, is_superuser, email, created_at, updated_at
	) VALUES (
	  @userName, @isSuperUser, @email, @now, @now
	)
	`
	bindParams := map[string]any{
		"userName":    user.UserName,
		"isSuperUser": user.IsSuperUser,
		"email":       user.Email,
		"now":         time.Now(),
	}
	if err := ur.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}

func (ur *userRepository) Update(ctx context.Context, user model.User) error {
	query := `
	UPDATE user
	SET
	  user_name = @userName,
	  is_superuser = @isSuperUser,
	  email = @email,
	  updated_at = @updatedAt,
	  updated_by = @userID
	WHERE user_id = @userID
	`
	bindParams := map[string]any{
		"userID":      user.UserID,
		"userName":    user.UserName,
		"isSuperUser": user.IsSuperUser,
		"email":       user.Email,
		"updatedAt":   time.Now(),
	}
	if err := ur.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}
