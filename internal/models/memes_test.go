package models

import (
	"testing"

	"github.com/disintegration/imaging"
)

// Image is > desired width
func Test_Thumnail(t *testing.T) {
	const testImgSrcPath = "fixtures/doom_meow.jpg"
	const expectedWidth = 200
	src, err := imaging.Open(testImgSrcPath)
	if err != nil {
		t.Errorf("Failed to open fixture image: %v", err)
	}

	actual := makeThumb(src, expectedWidth)
	actualWidth := actual.Rect.Dx()
	if actualWidth != expectedWidth {
		t.Errorf("Expected thumbnail to be %d pixels wide but got %d", expectedWidth, actualWidth)
	}
}

// Image is < desired width; should leave image alone
