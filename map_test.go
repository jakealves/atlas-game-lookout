package main_test

import (
	"image/png"
	"os"
	"testing"

	lookout "github.com/jakealves/atlas-game-lookout"
)

func TestConvertLongLatToPixel(t *testing.T) {
	tileX, tileY, pixelX, pixelY := lookout.ConvertCoordinatesToPixels(32, 58.52, 37.91)
	if tileX != 25 {
		t.Errorf("Expected tileX to be 25 got %v", tileX)
	}
	if tileY != 9 {
		t.Errorf("Expected tileY to be 9 got %v", tileY)
	}
	if pixelX != 93 {
		t.Errorf("Expected pixelX to be 93 got %v", pixelX)
	}
	if pixelY != 239 {
		t.Errorf("Expected pixelY to be 239 got %v", pixelY)
	}
}

func TestLoadImageTiles(t *testing.T) {
	test := []string{
		"test/50_19.png",
		"test/50_20.png",
		"test/51_19.png",
		"test/51_20.png",
	}
	ImageTiles, err := lookout.LoadImageTiles(test)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if len(ImageTiles) != 4 {
		t.Errorf("Expected splice to have a length of 4, instead it's %v", len(ImageTiles))
	}
	testImage, err := lookout.CombineImageTiles(ImageTiles, 2)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if testImage.Bounds().Max.X != 512 {
		t.Errorf("Expected testImage to have a width of 512, got %v", testImage.Bounds().Max.X)
	}
	if testImage.Bounds().Max.Y != 512 {
		t.Errorf("Expected testImage to have a width of 512, got %v", testImage.Bounds().Max.Y)
	}
	output, err := os.Create("testing.png")
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	defer output.Close()

	err = png.Encode(output, testImage)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
}

func TestGenerateTileFromCoordinates(t *testing.T) {
	imagePath, err := lookout.GenerateTileFromCoordinates("H3: Peg-Legged Peach demolished a 'Storage Box (Pin Coded)'", 64, 58.52, 37.91, 0)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if imagePath != "58_37.png" {
		t.Errorf("Expected imagePath to be 58_37.png got %v", imagePath)
	}

	imagePath2, err := lookout.GenerateTileFromCoordinates("", 64, -88.51, 53.12, 1)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if imagePath2 != "-88_53.png" {
		t.Errorf("Expected imagePath to be -88_53.png got %v", imagePath2)
	}
}
