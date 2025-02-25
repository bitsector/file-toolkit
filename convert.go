package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

// convertJPG converts the decoded image into a JPEG file at outPath.
func convertJPG(img image.Image, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("unable to create jpg file: %w", err)
	}
	defer outFile.Close()

	// Encode the image as JPEG with quality 100.
	if err := jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100}); err != nil {
		return fmt.Errorf("failed to encode jpg: %w", err)
	}
	return nil
}

// convertPNG converts the decoded image into a PNG file at outPath.
func convertPNG(img image.Image, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("unable to create png file: %w", err)
	}
	defer outFile.Close()

	// Encode the image as PNG.
	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}
	return nil
}

// logFileSizes prints the file sizes (in bytes) of all provided file paths.
func logFileSizes(paths ...string) error {
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("failed to get file stats for %s: %w", p, err)
		}
		fmt.Printf("%s size: %d bytes\n", filepath.Base(p), info.Size())
	}
	return nil
}

// confvert reads a .webp file, decodes it, and converts it to both JPEG and PNG formats.
// It then logs the file sizes for the original and converted images.
func confvert(filePath string) error {
	// Ensure the file has a ".webp" extension.
	if strings.ToLower(filepath.Ext(filePath)) != ".webp" {
		return fmt.Errorf("provided file is not a .webp file")
	}

	// Open the WebP file.
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	// Decode the WebP image.
	img, err := webp.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode webp image: %w", err)
	}

	// Generate output file names for JPEG and PNG.
	base := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	jpgFile := base + ".jpg"
	pngFile := base + ".png"

	// Convert and save the image as JPEG.
	if err := convertJPG(img, jpgFile); err != nil {
		return err
	}
	// Convert and save the image as PNG.
	if err := convertPNG(img, pngFile); err != nil {
		return err
	}

	// Log file sizes for the original WebP, JPEG, and PNG.
	if err := logFileSizes(filePath, jpgFile, pngFile); err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <image.webp>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	if err := confvert(filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
