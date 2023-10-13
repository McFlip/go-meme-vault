package tags

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Test_GetAll(t *testing.T) {
	const path = ":memory:"

	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&Tag{})
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
	actual := tagsModel.GetAll()

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
