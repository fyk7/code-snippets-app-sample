package generated

import "github.com/fyk7/code-snippets-app/app/domain/model"

// ToSnippetResponse converts a domain Snippet to a generated Snippet DTO.
func ToSnippetResponse(s model.Snippet) Snippet {
	resp := Snippet{
		SnipetId:           s.SnippetID,
		Title:              s.Title,
		Body:               s.Body,
		ProgramingLanguage: s.ProgramingLanguage,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
	}
	if s.Description != "" {
		resp.Description = &s.Description
	}
	if s.CreatedBy != 0 {
		resp.CreatedBy = &s.CreatedBy
	}
	if s.UpdatedBy != 0 {
		resp.UpdatedBy = &s.UpdatedBy
	}
	return resp
}

// ToSnippetResponses converts a slice of domain Snippets.
func ToSnippetResponses(snippets []model.Snippet) []Snippet {
	result := make([]Snippet, len(snippets))
	for i, s := range snippets {
		result[i] = ToSnippetResponse(s)
	}
	return result
}

// ToTagResponse converts a domain Tag to a generated Tag DTO.
func ToTagResponse(t model.Tag) Tag {
	resp := Tag{
		TagId:     t.TagID,
		TagName:   t.TagName,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.CreatedBy != 0 {
		resp.CreatedBy = &t.CreatedBy
	}
	if t.UpdatedBy != 0 {
		resp.UpdatedBy = &t.UpdatedBy
	}
	return resp
}

// ToTagResponses converts a slice of domain Tags.
func ToTagResponses(tags []model.Tag) []Tag {
	result := make([]Tag, len(tags))
	for i, t := range tags {
		result[i] = ToTagResponse(t)
	}
	return result
}

// SnippetCreateToModel converts a generated SnippetCreateRequest to a domain model.
func SnippetCreateToModel(req SnippetCreateRequest) model.Snippet {
	s := model.Snippet{
		Title:              req.Title,
		Body:               req.Body,
		ProgramingLanguage: req.ProgramingLanguage,
	}
	if req.Description != nil {
		s.Description = *req.Description
	}
	return s
}

// SnippetUpdateToModel converts a generated SnippetUpdateRequest to a domain model.
func SnippetUpdateToModel(req SnippetUpdateRequest) model.Snippet {
	s := model.Snippet{
		SnippetID:          uint64(req.SnippetId),
		Title:              req.Title,
		Body:               req.Body,
		ProgramingLanguage: req.ProgramingLanguage,
	}
	if req.Description != nil {
		s.Description = *req.Description
	}
	return s
}

// TagCreateToModel converts a generated TagCreateRequest to a domain model.
func TagCreateToModel(req TagCreateRequest) model.Tag {
	return model.Tag{
		TagName: req.TagName,
	}
}
