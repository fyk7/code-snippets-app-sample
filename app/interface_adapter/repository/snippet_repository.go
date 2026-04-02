package interface_adapter

import (
	"context"
	"strconv"
	"time"

	"github.com/fyk7/code-snippets-app/app/domain/model"
	"github.com/fyk7/code-snippets-app/app/domain/repository"
	"gorm.io/gorm"
)

type snippetRepository struct {
	Conn *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) repository.SnippetRepository {
	return &snippetRepository{
		Conn: db,
	}
}

func (sr *snippetRepository) GetAll(ctx context.Context) ([]model.Snippet, error) {
	var snippets []model.Snippet
	query := `SELECT * FROM snippet`
	if err := sr.Conn.WithContext(ctx).Raw(query).Scan(&snippets).Error; err != nil {
		return nil, toDomainError(err)
	}
	return snippets, nil
}

func (sr *snippetRepository) GetByID(ctx context.Context, snippetID uint64) (model.Snippet, error) {
	var snippet model.Snippet
	query := `
	SELECT *
	FROM snippet
	WHERE snippet_id = @snippetID
	`
	bindParams := map[string]any{
		"snippetID": snippetID,
	}
	if err := sr.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&snippet).Error; err != nil {
		return model.Snippet{}, toDomainError(err)
	}
	return snippet, nil
}

func (sr *snippetRepository) FindByKeyWord(ctx context.Context, keyword string) ([]model.Snippet, error) {
	var snippets []model.Snippet
	query := `
	SELECT *
	FROM snippet
	WHERE title LIKE CONCAT('%', @keyword, '%')
	   OR description LIKE CONCAT('%', @keyword, '%')
	`
	bindParams := map[string]any{
		"keyword": keyword,
	}
	if err := sr.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&snippets).Error; err != nil {
		return nil, toDomainError(err)
	}
	return snippets, nil
}

func (sr *snippetRepository) FindByTag(ctx context.Context, tagID uint64) ([]model.Snippet, error) {
	var snippets []model.Snippet
	query := `
	SELECT s.*
	FROM snippet s
	INNER JOIN snippet_tag_relation r
	  ON s.snippet_id = r.snippet_id
	WHERE r.tag_id = @tagID
	`
	bindParams := map[string]any{
		"tagID": tagID,
	}
	if err := sr.Conn.WithContext(ctx).Raw(query, bindParams).Scan(&snippets).Error; err != nil {
		return nil, toDomainError(err)
	}
	return snippets, nil
}

func (sr *snippetRepository) AssociateWithTag(ctx context.Context, snippetID, tagID, userID int64) error {
	query := `
	INSERT INTO snippet_tag_relation (
	  snippet_id, tag_id, created_at, created_by, updated_at, updated_by
	) VALUES (
	  @snippetID, @tagID, @now, @userID, @now, @userID
	)
	`
	now := time.Now()
	bindParams := map[string]any{
		"snippetID": snippetID,
		"tagID":     tagID,
		"now":       now,
		"userID":    strconv.FormatInt(userID, 10),
	}
	if err := sr.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}

func (sr *snippetRepository) Create(ctx context.Context, snippet model.Snippet, userID uint64) error {
	query := `
	INSERT INTO snippet (
	  title, description, body, programing_language,
	  created_at, created_by, updated_at, updated_by
	) VALUES (
	  @title, @description, @body, @programingLanguage,
	  @now, @userID, @now, @userID
	)
	`
	now := time.Now()
	bindParams := map[string]any{
		"title":              snippet.Title,
		"description":        snippet.Description,
		"body":               snippet.Body,
		"programingLanguage": snippet.ProgramingLanguage,
		"now":                now,
		"userID":             strconv.FormatUint(userID, 10),
	}
	if err := sr.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}

func (sr *snippetRepository) Update(ctx context.Context, snippet model.Snippet, userID uint64) error {
	query := `
	UPDATE snippet
	SET
	  title = @title,
	  description = @description,
	  body = @body,
	  programing_language = @programingLanguage,
	  updated_at = @now,
	  updated_by = @userID
	WHERE snippet_id = @id
	`
	bindParams := map[string]any{
		"id":                 snippet.SnippetID,
		"title":              snippet.Title,
		"description":        snippet.Description,
		"body":               snippet.Body,
		"programingLanguage": snippet.ProgramingLanguage,
		"now":                time.Now(),
		"userID":             strconv.FormatUint(userID, 10),
	}
	if err := sr.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}

func (sr *snippetRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM snippet WHERE snippet_id = @id`
	bindParams := map[string]any{
		"id": id,
	}
	if err := sr.Conn.WithContext(ctx).Exec(query, bindParams).Error; err != nil {
		return toDomainError(err)
	}
	return nil
}
