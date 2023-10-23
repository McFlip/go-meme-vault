package models

import (
	"image"

	"github.com/disintegration/imaging"
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

func makeThumb(src image.Image, width int) image.Image {
	if src.Bounds().Dx() > width {
		return imaging.Resize(src, width, 0, imaging.Box)
	}
	return src
}
