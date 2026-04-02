package interface_adapter

import "github.com/fyk7/code-snippets-app/app/domain/model"

type SnippetPostReq struct {
	Title              string `json:"title" validate:"required"`
	Description        string `json:"description"`
	Body               string `json:"body" validate:"required"`
	ProgramingLanguage string `json:"programing_language" validate:"required"`
}

func (spr *SnippetPostReq) ConvertToModel() model.Snippet {
	return model.Snippet{
		Title:              spr.Title,
		Description:        spr.Description,
		Body:               spr.Body,
		ProgramingLanguage: spr.ProgramingLanguage,
	}
}

type SnippetPutReq struct {
	SnippetID          int64  `json:"snippet_id" validate:"required"`
	Title              string `json:"title" validate:"required"`
	Description        string `json:"description"`
	Body               string `json:"body" validate:"required"`
	ProgramingLanguage string `json:"programing_language" validate:"required"`
}

func (spr *SnippetPutReq) ConvertToModel() model.Snippet {
	return model.Snippet{
		SnippetID:          uint64(spr.SnippetID),
		Title:              spr.Title,
		Description:        spr.Description,
		Body:               spr.Body,
		ProgramingLanguage: spr.ProgramingLanguage,
	}
}
