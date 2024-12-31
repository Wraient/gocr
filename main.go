package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func printHelp() {
	fmt.Println("OCR Tool - Extract text from images")
	fmt.Println("\nUsage:")
	fmt.Println("  gocr [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -i, --image string    Process image file in CLI mode")
	fmt.Println("  -g, --gui string      Open GUI with specified image")
	fmt.Println("  -h, --help            Show help message")
	fmt.Println("\nExamples:")
	fmt.Println("  gocr                  Launch GUI application")
	fmt.Println("  gocr -i image.png     Process image in CLI mode")
	fmt.Println("  gocr -g image.png     Open GUI with image loaded")
}

func main() {
	// Define flags
	imagePath := flag.String("i", "", "Path to image file for OCR processing")
	guiImage := flag.String("g", "", "Open GUI with specified image")
	help := flag.Bool("h", false, "Show help message")

	// Parse flags
	flag.Parse()

	// Show help if requested
	if *help {
		printHelp()
		os.Exit(0)
	}

	// Create an instance of the app structure
	app := NewApp()

	// If image path is provided for CLI processing
	if *imagePath != "" {
		result, err := app.ProcessImageFile(*imagePath)
		if err != nil {
			fmt.Printf("Error processing image: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result.Text)
		os.Exit(0)
	}

	// If GUI image path is provided, set it as initial image
	if *guiImage != "" {
		absPath, err := filepath.Abs(*guiImage)
		if err != nil {
			fmt.Printf("Error resolving path: %v\n", err)
			os.Exit(1)
		}
		app.InitialImage = absPath
	}

	// Launch GUI
	err := wails.Run(&options.App{
		Title:  "OCR Tool",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
