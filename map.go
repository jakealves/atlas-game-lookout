package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func ConvertCoordinatesToPixels(maxDimension int, long float64, lat float64) (tileX, tileY, pixelX, pixelY int) {
	// get total number of pixels all the tiles equal to
	totalPixels := maxDimension * 256

	// determine how many pixels are in a single lat or long
	oneDegree := float64(totalPixels) / float64(200)

	var longAsPixels float64
	var latAsPixels float64

	// check if long is negative
	if math.Signbit(long) {
		// how many total pixels from the left are we?
		longAsPixels = (100 - math.Abs(long)) * oneDegree
	} else {
		// how many total pixels from the left are we? plus add 100
		// cause a positive long means we are over halfway there.
		longAsPixels = (float64(100) + long) * oneDegree
	}

	// check if lat is negative
	if math.Signbit(lat) {
		// how many total pixels from the top are we? plus add 100
		// cause a negative means we are over halfway there.
		latAsPixels = (float64(100) + math.Abs(lat)) * oneDegree
	} else {
		// how many total pixels from the top are we? minus 100 because
		// lat 100 is actually pixel 0.
		latAsPixels = (float64(100) - lat) * oneDegree
	}

	// determine what x tile to use
	xRaw := longAsPixels / float64(256)
	tileX = int(math.Floor(xRaw))

	// use the remainder as a percent to determine the x pixel
	pixelX = int(math.Round((xRaw - float64(tileX)) * float64(256)))

	// determine what y tile to use
	yRaw := latAsPixels / float64(256)
	tileY = int(math.Floor(yRaw))

	// use the remainder as a percent to determine the y pixel
	pixelY = int(math.Round((yRaw - float64(tileY)) * float64(256)))
	return tileX, tileY, pixelX, pixelY
}

func LoadImageTiles(tileList []string) ([]image.Image, error) {
	var imageList []image.Image
	for _, fileName := range tileList {
		imageFile, err := os.Open(fileName)
		if err != nil {
			return nil, fmt.Errorf("%s -> %s", err, fileName)
		}
		defer imageFile.Close()
		pngImage, err := png.Decode(imageFile)
		if err != nil {
			return nil, err
		}
		imageList = append(imageList, pngImage)
	}
	return imageList, nil
}

func CombineImageTiles(imageList []image.Image, gridSize int) (*image.RGBA, error) {
	totalTiles := len(imageList)
	newImageDimensions := (totalTiles / gridSize) * 256

	newImage := image.NewRGBA(image.Rect(0, 0, newImageDimensions, newImageDimensions))
	var bounds image.Rectangle
	initX := 0
	initY := 0

	// there has got to be a better way
	for i, img := range imageList {
		if i == gridSize {
			initX = initX + 256
			initY = 0
		}
		bounds = image.Rect(initX, initY, newImageDimensions, newImageDimensions)
		initY = initY + 256
		draw.Draw(newImage, bounds, img, image.Point{}, draw.Src)
	}

	return newImage, nil
}

func ShowSurroundingTiles(initialTile string, initialIndex, maxDimension, startX, startY int) (*image.RGBA, error) {
	var tileList []string
	for i := 0; i < 4; i++ {
		var filePath string
		switch index := i; index {
		case 0:
			filePath = fmt.Sprintf("tiles/%d/%d_%d.png", maxDimension, startX, startY)
		case 1:
			filePath = fmt.Sprintf("tiles/%d/%d_%d.png", maxDimension, startX, startY+1)
		case 2:
			filePath = fmt.Sprintf("tiles/%d/%d_%d.png", maxDimension, startX+1, startY)
		case 3:
			filePath = fmt.Sprintf("tiles/%d/%d_%d.png", maxDimension, startX+1, startY+1)
		}
		tileList = append(tileList, filePath)
	}

	tileList[initialIndex] = initialTile

	imageList, err := LoadImageTiles(tileList)
	if err != nil {
		return nil, fmt.Errorf("there was an error executing LoadImageTiles-> %s", err)
	}

	newImage, err := CombineImageTiles(imageList, 2)
	if err != nil {
		return nil, fmt.Errorf("there was an error executing CombineImageTiles-> %s", err)

	}
	return newImage, nil
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}

	point := fixed.Point26_6{}
	point.X = fixed.I(x)
	point.Y = fixed.I(y)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func AddRedDotUsingCoordinates(inputFile string, xCoordinate int, yCoordinate int, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode input png: %w", err)
	}

	bounds := img.Bounds()

	// Create a new image with the same dimensions as the input image
	newImg := image.NewRGBA(bounds)

	// Copy the input image to the new image
	draw.Draw(newImg, bounds, img, image.Point{}, draw.Src)

	// Add a red dot to the center of the new image
	redDot := color.RGBA{255, 0, 0, 255}
	drawCircle(newImg, xCoordinate, yCoordinate, 4, redDot)

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	if err := png.Encode(output, newImg); err != nil {
		return fmt.Errorf("failed to encode output png: %w", err)
	}

	return nil
}

func drawCircle(img *image.RGBA, centerX, centerY, radius int, circleColor color.Color) {
	sqRadius := radius * radius

	// Iterate through the pixels within a square bounding the circle
	for x := centerX - radius; x <= centerX+radius; x++ {
		for y := centerY - radius; y <= centerY+radius; y++ {
			sqDist := (x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)

			// Check if the pixel is inside the circle using the distance formula
			if sqDist <= sqRadius {
				img.Set(x, y, circleColor)
			}
		}
	}
}

func DownloadMapTiles(baseurl string, zoom int, max int) error {
	for x := 0; x <= max; x++ {
		for y := 0; y <= max; y++ {
			tileURL := fmt.Sprintf("%s/%s/%s/%s.png", baseurl, fmt.Sprint(zoom), fmt.Sprint(x), fmt.Sprint(y))
			fileName := fmt.Sprintf("tiles/%s/%s_%s.png", fmt.Sprint(zoom), fmt.Sprint(x), fmt.Sprint(y))
			DownloadFile(tileURL, fileName)
		}
	}
	return nil
}

func DownloadFile(URL, fileName string) error {

	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func UnzipFile(fileName string) error {
	archive, err := zip.OpenReader(fileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := f.Name
		fmt.Println("unzipping file ", filePath)

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

func GenerateTileFromCoordinates(label string, maxDimension int, long float64, lat float64, shiftIndex int) (imagePath string, err error) {
	tileX, tileY, pixelX, pixelY := ConvertCoordinatesToPixels(maxDimension, long, lat)

	sourceTile := fmt.Sprintf("tiles/%d/%d_%d.png", maxDimension, tileX, tileY)
	outputTile := fmt.Sprintf("%d_%d.png", int(long), int(lat))

	err = AddRedDotUsingCoordinates(sourceTile, pixelX, pixelY, outputTile)
	if err != nil {
		return "", fmt.Errorf("unexpected error drawing red dot on image -> %v", err)
	}

	var image *image.RGBA
	switch initialIndex := shiftIndex; initialIndex {
	case 0:
		image, err = ShowSurroundingTiles(outputTile, initialIndex, maxDimension, tileX, tileY)
	case 1:
		image, err = ShowSurroundingTiles(outputTile, initialIndex, maxDimension, tileX, tileY-1)
	case 2:
		image, err = ShowSurroundingTiles(outputTile, initialIndex, maxDimension, tileX-1, tileY)
	case 3:
		image, err = ShowSurroundingTiles(outputTile, initialIndex, maxDimension, tileX-1, tileY-1)
	}

	addLabel(image, 25, 25, label)
	addLabel(image, 25, 487, fmt.Sprintf("Long:%.2f Lat:%.2f", long, lat))

	output, err := os.Create(outputTile)
	if err != nil {
		return "", fmt.Errorf("there was an error creating %s -> %s", outputTile, err)
	}
	defer output.Close()

	if err := png.Encode(output, image); err != nil {
		return "", fmt.Errorf("error encoding new image-> %s", err)
	}

	if err != nil {
		return "", fmt.Errorf("unexpected error stiching surrounding tiles -> %v", err)
	}
	return outputTile, nil
}
