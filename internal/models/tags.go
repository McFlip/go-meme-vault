package models

import (
	"gorm.io/gorm"
)

type TagsModel struct {
	DB *gorm.DB
}

type Tag struct {
	gorm.Model
	Name  string
	Memes []*Meme `gorm:"many2many:meme_tags;"`
}

func (tagsModel *TagsModel) GetAll() ([]Tag, error) {
	var tags []Tag
	res := tagsModel.DB.Find(&tags)
	return tags, res.Error
}

func (tagsModel *TagsModel) GetByID(id uint) (Tag, error) {
	var tag Tag
	res := tagsModel.DB.First(&tag, id)
	return tag, res.Error
}

func (tagsModel *TagsModel) Create(tag *Tag) *gorm.DB {
	newTag := tagsModel.DB.Create(tag)
	return newTag
}
