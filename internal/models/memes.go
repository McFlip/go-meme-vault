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

func (memesModel *MemesModel) Create(meme *Meme) *gorm.DB {
	return memesModel.DB.Create(meme)
}

func (memesModel *MemesModel) GetByID(id uint) (Meme, error) {
	var meme Meme
	res := memesModel.DB.First(&meme, id)
	return meme, res.Error
}

func (memesModel *MemesModel) GetByPath(path string) (Meme, error) {
	var meme Meme
	res := memesModel.DB.Where(&Meme{Path: path}).First(&meme)
	return meme, res.Error
}
