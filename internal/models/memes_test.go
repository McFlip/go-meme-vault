package models

import (
	"testing"

	"github.com/disintegration/imaging"
)

// Image is > desired width
func Test_Thumbnail(t *testing.T) {
	const testImgSrcPath = "fixtures/doom_meow.jpg"
	const expectedWidth = 200
	src, err := imaging.Open(testImgSrcPath)
	if err != nil {
		t.Errorf("Failed to open fixture image: %v", err)
	}

	actual := makeThumb(src, expectedWidth)
	actualWidth := actual.Bounds().Dx()
	if actualWidth != expectedWidth {
		t.Errorf("Expected thumbnail to be %d pixels wide but got %d", expectedWidth, actualWidth)
	}
}

// Image is <= desired width; should leave image alone
func Test_ThumbnailTooSmall(t *testing.T) {
	const testImgSrcPath = "fixtures/doom_meow_tn.jpg" // 200 x 233
	const expectedWidth = 200
	src, err := imaging.Open(testImgSrcPath)
	if err != nil {
		t.Errorf("Failed to open fixture image: %v", err)
	}

	// We ask to "downscale" to a width > src width, so no scaling should happen
	actual := makeThumb(src, 300)
	actualWidth := actual.Bounds().Dx()
	if actualWidth != expectedWidth {
		t.Errorf("Expected thumbnail to be %d pixels wide but got %d", expectedWidth, actualWidth)
	}
}

func Test_CreateMeme(t *testing.T) {
	const expectedPath = "fixtures/doom_meow.jpg"
	testMeme := Meme{
		Path: expectedPath,
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}

	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create meme: %s", res.Error)
	}

	actual, err := testModel.GetByID(testMeme.ID)
	if res.Error != nil {
		t.Errorf("Failed to get meme: %s", err)
	}

	if actual.Path != expectedPath {
		t.Errorf("Expected new meme to have path %s, but got %s", expectedPath, actual.Path)
	}
}

func Test_GetByPath(t *testing.T) {
	const expectedPath = "fixtures/doom_meow.jpg"
	testMeme := Meme{
		Path: expectedPath,
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}

	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create meme: %s", res.Error)
	}

	actual, err := testModel.GetByPath(expectedPath)
	if res.Error != nil {
		t.Errorf("Failed to get meme: %s", err)
	}

	if actual.Path != expectedPath {
		t.Errorf("Expected new meme to have path %s, but got %s", expectedPath, actual.Path)
	}
}

func Test_GetAllMemes(t *testing.T) {
	const expectedPath = "fixtures/doom_meow.jpg"
	testMeme := Meme{
		Path: expectedPath,
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}

	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create meme: %s", res.Error)
	}

	actual, err := testModel.GetAll()
	if res.Error != nil {
		t.Errorf("Failed to get meme: %s", err)
	}

	if len(actual) != 1 {
		t.Errorf("Expected memes count to be 1, but got %d", len(actual))
	} else if actual[0].Path != expectedPath {
		t.Errorf("Expected new meme to have path %s, but got %s", expectedPath, actual[0].Path)
	}
}
