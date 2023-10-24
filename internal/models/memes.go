package models

import (
	"errors"
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

const Tn_Width = 200

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

func (m *Meme) GetThumbnail() string {
	return strings.Replace(m.Path, "full", "tn", 1)
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

func (memesModel *MemesModel) Scan(path string) ([]Meme, error) {
	var memes []Meme
	// walk path
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// check if img already exists as meme
		if !d.IsDir() {
			_, err := memesModel.GetByPath(d.Name())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// create the meme and thumb
				newMeme := Meme{
					Name: filepath.Base(path),
					Path: fmt.Sprintf("/%s", path),
				}
				res := memesModel.Create(&newMeme)
				if res.Error != nil {
					return res.Error
				}
				src, err := imaging.Open(path)
				if err != nil {
					return err
				}
				thumbnail := makeThumb(src, Tn_Width)
				imaging.Save(thumbnail, newMeme.GetThumbnail())
				// append new meme to list
				memes = append(memes, newMeme)
			} else if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return memes, err
	}
	return memes, nil
}
