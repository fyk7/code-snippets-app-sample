package interface_adapter

import "github.com/fyk7/code-snippets-app/app/domain/model"

type TagPostReq struct {
	TagName string `json:"tag_name" validate:"required"`
}

func (spr *TagPostReq) ConvertToModel() model.Tag {
	return model.Tag{
		TagName: spr.TagName,
	}
}

type TagPutReq struct {
	TagID   uint64 `json:"tag_id" validate:"required"`
	TagName string `json:"tag_name" validate:"required"`
}

func (spr *TagPutReq) ConvertToModel() model.Tag {
	return model.Tag{
		TagID:   spr.TagID,
		TagName: spr.TagName,
	}
}
