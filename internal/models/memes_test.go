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

func Test_FilterTags(t *testing.T) {
	testTag1 := Tag{
		Name: "First Tag",
	}
	testTag2 := Tag{
		Name: "Second Tag",
	}
	testMeme := Meme{
		Tags: []*Tag{&testTag1},
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}
	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create test meme: %s", res.Error)
	}

	actual, err := testModel.FilterTags(testMeme.ID, []Tag{testTag1, testTag2})
	if err != nil {
		t.Errorf("Failed to filter tags: %s", err)
	}

	if len(actual) != 1 {
		t.Errorf("Expected tags length of 1, but got %v", len(actual))
	} else {
		if actual[0].Name != testTag2.Name {
			t.Errorf("Expected filtered tag name to be %s, but got %s", testTag2.Name, actual[0].Name)
		}
	}
}

func Test_GetUntagged(t *testing.T) {
	testTag := Tag{
		Name: "Test Tag",
	}
	testMemeWithTag := Meme{
		Name: "With Tag",
		Tags: []*Tag{&testTag},
	}
	testMemeNoTag := Meme{
		Name: "no tag",
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}
	res := testModel.Create(&testMemeWithTag)
	if res.Error != nil {
		t.Errorf("Failed to create test meme with tag: %s", res.Error)
	}
	res = testModel.Create(&testMemeNoTag)
	if res.Error != nil {
		t.Errorf("Failed to create test meme without tag: %s", res.Error)
	}

	actual, err := testModel.GetUntagged()
	if err != nil {
		t.Errorf("Failted to get untagged memes: %s", err)
	}

	if len(actual) != 1 {
		t.Errorf("Expected length of meme slice to be 1, but got %d", len(actual))
	} else {
		if actual[0].Name != "no tag" {
			t.Errorf("Expected tag to have name 'no tag', but got %v", actual)
		}
	}
}

func Test_AddTag(t *testing.T) {
	testTag1 := Tag{
		Name: "First Tag",
	}
	testTag2 := Tag{
		Name: "Second Tag",
	}
	testMeme := Meme{
		Tags: []*Tag{&testTag1},
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}
	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create test meme: %s", res.Error)
	}

	actual, err := testModel.AddTag(testMeme.ID, testTag2)
	if err != nil {
		t.Errorf("Failed to add tag: %s", err)
	}

	if len(actual.Tags) != 2 {
		t.Errorf("Expected tags length of 2, but got %v", len(actual.Tags))
	} else {
		if actual.Tags[1].Name != testTag2.Name {
			t.Errorf("Expected 2nd tag name to be %s, but got %s", testTag2.Name, actual.Tags[1].Name)
		}
	}
}

func Test_RemoveTag(t *testing.T) {
	testTag1 := Tag{
		Name: "First Tag",
	}
	testMeme := Meme{
		Tags: []*Tag{&testTag1},
	}
	database := connectDB()
	testModel := MemesModel{
		DB: database,
	}
	res := testModel.Create(&testMeme)
	if res.Error != nil {
		t.Errorf("Failed to create test meme: %s", res.Error)
	}

	err := testModel.RemoveTag(testMeme.ID, testTag1)
	if err != nil {
		t.Errorf("Failed to remove tag: %s", err)
	}

	actual, err := testModel.GetByID(testMeme.ID)
	if err != nil {
		t.Errorf("Failed to get test meme: %s", err)
	}
	if len(actual.Tags) != 0 {
		t.Errorf("Expected tags length of 0, but got %v", len(actual.Tags))
	}
}
