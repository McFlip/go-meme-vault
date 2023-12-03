package models

import (
	"errors"
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"slices"
	"sort"
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
	res := memesModel.DB.Preload("Tags").First(&meme, id)
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
			_, err := memesModel.GetByPath(fmt.Sprintf("/%s", path))
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

func (memesModel *MemesModel) GetAll() ([]Meme, error) {
	var memes []Meme
	res := memesModel.DB.Find(&memes)
	return memes, res.Error
}

func (memesModel *MemesModel) GetUntagged() ([]Meme, error) {
	var memes []Meme
	res := memesModel.DB.Raw("select * from memes where id not in (select distinct meme_id from meme_tags)").Scan(&memes)
	return memes, res.Error
}

func (memesModel *MemesModel) FilterTags(id uint, tags []Tag) ([]Tag, error) {
	meme, err := memesModel.GetByID(id)
	if err != nil {
		return nil, err
	}

	filteredTags := make([]Tag, 0, len(tags))
	tagIds := make(sort.IntSlice, 0, len(meme.Tags))
	for _, t := range meme.Tags {
		tagIds = append(tagIds, int(t.ID))
	}
	tagIds.Sort()
	for _, t := range tags {
		_, found := slices.BinarySearch(tagIds, int(t.ID))
		if !found {
			filteredTags = append(filteredTags, t)
		}
	}

	return filteredTags, nil
}

func (memesModel *MemesModel) AddTag(id uint, tag Tag) (Meme, error) {
	meme, err := memesModel.GetByID(id)
	if err != nil {
		return meme, err
	}

	err = memesModel.DB.Model(&meme).Association("Tags").Append(&tag)

	return meme, err
}

func (memesModel *MemesModel) RemoveTag(id uint, tag Tag) error {
	meme, err := memesModel.GetByID(id)
	if err != nil {
		return err
	}

	err = memesModel.DB.Model(&meme).Association("Tags").Delete(&tag)

	return err
}
