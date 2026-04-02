package interface_adapter

import (
	"context"
	"strconv"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/domain/repository"
	"gorm.io/gorm"
)

type tagRepository struct {
	Conn *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{
		Conn: db,
	}
}

func (tr *tagRepository) GetAll(ctx context.Context) ([]model.Tag, error) {
	var tags []model.Tag
	query := `SELECT * FROM tag`
	if err := tr.Conn.WithContext(ctx).Raw(query).Scan(&tags).Error; err != nil {
		return nil, toDomainError(err)
	}
	return tags, nil
}

func (tr *tagRepository) GetByID(ctx context.Context, tagID uint64) (model.Tag, error) {
	var tag model.Tag
	query := `
	SELECT *
	FROM tag
	WHERE tag_id = @tagID
	`
	bindParams := map[string]any{
		"tagID": tagID,
	}
	if err := tr.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&tag).Error; err != nil {
		return model.Tag{}, toDomainError(err)
	}
	return tag, nil
}

func (tr *tagRepository) FindByKeyWord(ctx context.Context, keyword string) ([]model.Tag, error) {
	var tags []model.Tag
	query := `
	SELECT *
	FROM tag
	WHERE tag_name LIKE CONCAT('%', @keyword, '%')
	`
	bindParams := map[string]any{
		"keyword": keyword,
	}
	if err := tr.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&tags).Error; err != nil {
		return nil, toDomainError(err)
	}
	return tags, nil
}

func (tr *tagRepository) Create(ctx context.Context, tag model.Tag, userID uint64) error {
	query := `
	INSERT INTO tag (
	  tag_name, created_at, created_by, updated_at, updated_by
	) VALUES (
	  @tagName, @now, @userID, @now, @userID
	)
	`
	now := time.Now()
	bindParams := map[string]any{
		"tagName": tag.TagName,
		"now":     now,
		"userID":  strconv.FormatUint(userID, 10),
	}
	if err := tr.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}
