package models

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	const path = ":memory:"

	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&Tag{})
	database.AutoMigrate(&Meme{})

	return database
}

func Test_GetAll(t *testing.T) {
	database := connectDB()
	expected := []Tag{
		{
			Name: "first tag",
		},
		{
			Name: "second tag",
		},
		{
			Name: "third tag",
		},
	}
	for _, t := range expected {
		database.Create(&t)
	}

	tagsModel := TagsModel{DB: database}
	actual, err := tagsModel.GetAll()
	if err != nil {
		t.Errorf("ERROR getting all tags: %s", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected %v but got %v", expected, actual)
	} else {
		for i, tag := range actual {
			if tag.Name != expected[i].Name {
				t.Errorf("Expected %v but got %v", expected[i].Name, actual[i].Name)
			}
		}
	}
}

func Test_Create(t *testing.T) {
	database := connectDB()
	expected := Tag{
		Name: "new hotness",
	}
	tagsmodel := TagsModel{
		DB: database,
	}

	res := tagsmodel.Create(&expected)
	if res.Error != nil {
		t.Errorf("ERROR in creating tag: %s", res.Error)
	}

	actual, err := tagsmodel.GetByID(expected.ID)
	if err != nil {
		t.Error(err)
	}

	if actual.Name != expected.Name {
		t.Errorf("Expected Tag Name to be %s, but got %s", expected.Name, actual.Name)
	}
}

func Test_TagExists(t *testing.T) {
	database := connectDB()
	expectedTrue := true
	expectedFalse := false
	trueName := "veritas"
	falseName := "cap"
	testTag := Tag{
		Name: trueName,
	}
	tagsmodel := TagsModel{
		DB: database,
	}

	res := tagsmodel.Create(&testTag)
	if res.Error != nil {
		t.Errorf("ERROR in creating tag: %s", res.Error)
	}

	if actual, ok := tagsmodel.TagExists(trueName); ok {
		if actual.Name != trueName {
			t.Errorf("Expected tag name to be %v, but got %v", trueName, actual.Name)
		}
	} else {
		t.Errorf("Expected tag exists to be %v, but got %v", expectedTrue, ok)
	}

	if _, ok := tagsmodel.TagExists(falseName); ok {
		t.Errorf("Expected tag exists to be %v, but got %v", expectedFalse, ok)
	}
}
