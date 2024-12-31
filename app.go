package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// App struct
type App struct {
	ctx context.Context
}

type OCRResult struct {
	Text  string `json:"text"`
	Boxes []Box  `json:"boxes"`
}

type Box struct {
	Text   string `json:"text"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// NewApp creates a new App application 
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ProcessImage(imageData string) (OCRResult, error) {
	// Remove data URL prefix if present
	imageData = strings.TrimPrefix(imageData, "data:image/jpeg;base64,")
	imageData = strings.TrimPrefix(imageData, "data:image/png;base64,")

	// Decode base64 image
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return OCRResult{}, err
	}

	// Initialize Tesseract client
	client := gosseract.NewClient()
	defer client.Close()

	// Set image from decoded bytes
	if err := client.SetImageFromBytes(decoded); err != nil {
		return OCRResult{}, err
	}

	// Get bounding boxes
	boxes, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
	if err != nil {
		return OCRResult{}, err
	}

	var result OCRResult
	for _, box := range boxes {
		result.Boxes = append(result.Boxes, Box{
			Text:   box.Word,
			X:      box.Box.Min.X,
			Y:      box.Box.Min.Y,
			Width:  box.Box.Max.X - box.Box.Min.X,
			Height: box.Box.Max.Y - box.Box.Min.Y,
		})
	}

	// Get full text
	text, err := client.Text()
	if err != nil {
		return OCRResult{}, err
	}
	result.Text = text

	return result, nil
}
