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
