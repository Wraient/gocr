package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// App struct
type App struct {
	ctx          context.Context
	InitialImage string
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

func (a *App) ProcessImage(input string) (OCRResult, error) {
	// If input starts with "data:", it's a base64 image
	if strings.HasPrefix(input, "data:") {
		// Extract the base64 data after the comma
		base64Data := strings.Split(input, ",")[1]
		imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return OCRResult{}, fmt.Errorf("error decoding base64 image: %v", err)
		}

		// Initialize Tesseract client
		client := gosseract.NewClient()
		defer client.Close()

		// Set OCR settings (same as ProcessImageFile)
		client.SetVariable("tessedit_pageseg_mode", "1")
		client.SetVariable("tessedit_ocr_engine_mode", "2")
		client.SetVariable("preserve_interword_spaces", "1")
		client.SetVariable("textord_heavy_nr", "1")
		client.SetVariable("textord_min_linesize", "2.5")
		client.SetVariable("tessedit_char_blacklist", "§¶©®™")

		if err := client.SetImageFromBytes(imageBytes); err != nil {
			return OCRResult{}, err
		}

		// Get bounding boxes and process OCR
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

		// Merge nearby boxes
		result.Boxes = mergeBoxes(result.Boxes)

		// Get full text
		text, err := client.Text()
		if err != nil {
			return OCRResult{}, err
		}
		result.Text = text

		return result, nil
	}

	return OCRResult{}, fmt.Errorf("invalid input format")
}

func (a *App) ProcessImageFile(filepath string) (OCRResult, error) {
	// Read the image file
	imageBytes, err := os.ReadFile(filepath)
	if err != nil {
		return OCRResult{}, fmt.Errorf("error reading image file: %v", err)
	}

	// Initialize Tesseract client
	client := gosseract.NewClient()
	defer client.Close()

	// Basic configuration
	if err := client.SetLanguage("eng"); err != nil {
		return OCRResult{}, err
	}

	// Optimize OCR settings
	client.SetVariable("tessedit_pageseg_mode", "1")       // Automatic page segmentation with OSD
	client.SetVariable("tessedit_ocr_engine_mode", "2")    // Legacy + LSTM mode
	client.SetVariable("preserve_interword_spaces", "1")   // Preserve spacing between words
	client.SetVariable("textord_heavy_nr", "1")            // Heavy noise removal
	client.SetVariable("textord_min_linesize", "2.5")      // Minimum text size to detect
	client.SetVariable("tessedit_char_blacklist", "§¶©®™") // Exclude problematic characters

	// Set image from bytes
	if err := client.SetImageFromBytes(imageBytes); err != nil {
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

	// Merge nearby boxes
	result.Boxes = mergeBoxes(result.Boxes)

	// Get full text
	text, err := client.Text()
	if err != nil {
		return OCRResult{}, err
	}
	result.Text = text

	return result, nil
}

func (b Box) isNearby(other Box) bool {
	// Consider boxes nearby if they are within 20 pixels horizontally
	// and 10 pixels vertically of each other
	horizontalGap := 20
	verticalGap := 10

	// Check if boxes are horizontally nearby
	horizontallyNear := (b.X <= other.X+other.Width+horizontalGap) &&
		(other.X <= b.X+b.Width+horizontalGap)

	// Check if boxes are vertically nearby
	verticallyNear := (b.Y <= other.Y+other.Height+verticalGap) &&
		(other.Y <= b.Y+b.Height+verticalGap)

	return horizontallyNear && verticallyNear
}

func mergeBoxes(boxes []Box) []Box {
	if len(boxes) == 0 {
		return boxes
	}

	// Sort boxes by Y coordinate first, then X coordinate
	sort.Slice(boxes, func(i, j int) bool {
		if abs(boxes[i].Y-boxes[j].Y) < 10 { // If Y coordinates are similar, sort by X
			return boxes[i].X < boxes[j].X
		}
		return boxes[i].Y < boxes[j].Y
	})

	var mergedBoxes []Box
	currentGroup := []Box{boxes[0]}

	for i := 1; i < len(boxes); i++ {
		// Check if current box is nearby any box in the current group
		isNearby := false
		for _, groupBox := range currentGroup {
			if groupBox.isNearby(boxes[i]) {
				isNearby = true
				break
			}
		}

		if isNearby {
			currentGroup = append(currentGroup, boxes[i])
		} else {
			// Merge current group and start a new one
			mergedBox := mergeGroup(currentGroup)
			mergedBoxes = append(mergedBoxes, mergedBox)
			currentGroup = []Box{boxes[i]}
		}
	}

	// Don't forget to merge the last group
	if len(currentGroup) > 0 {
		mergedBox := mergeGroup(currentGroup)
		mergedBoxes = append(mergedBoxes, mergedBox)
	}

	return mergedBoxes
}

func mergeGroup(group []Box) Box {
	if len(group) == 0 {
		return Box{}
	}
	if len(group) == 1 {
		return group[0]
	}

	// Find the bounding box that contains all boxes in the group
	minX := group[0].X
	minY := group[0].Y
	maxX := group[0].X + group[0].Width
	maxY := group[0].Y + group[0].Height
	words := []string{group[0].Text}

	for i := 1; i < len(group); i++ {
		box := group[i]
		words = append(words, box.Text)
		if box.X < minX {
			minX = box.X
		}
		if box.Y < minY {
			minY = box.Y
		}
		if box.X+box.Width > maxX {
			maxX = box.X + box.Width
		}
		if box.Y+box.Height > maxY {
			maxY = box.Y + box.Height
		}
	}

	return Box{
		Text:   strings.Join(words, " "),
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Add a new method to get the initial image
func (a *App) GetInitialImage() string {
	return a.InitialImage
}

func (a *App) GetImageData(input string) (string, error) {
	// If input is a file path, read the file and convert to base64
	imageBytes, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("error reading image file: %v", err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes), nil
}
