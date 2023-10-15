package models

import (
	"gorm.io/gorm"
)

type MemesModel struct {
	DB *gorm.DB
}

type Meme struct {
	gorm.Model
	Name string
	Path string
	Tags []*Tag `gorm:"many2many:meme_tags;"`
}
