package tags

import (
	"gorm.io/gorm"
)

type TagsModel struct {
	DB *gorm.DB
}

type Tag struct {
	gorm.Model
	Name string
}

func (tagsModel *TagsModel) GetAll() []Tag {
	var tags []Tag
	tagsModel.DB.Find(&tags)
	return tags
}