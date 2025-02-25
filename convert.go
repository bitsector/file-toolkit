package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

func confvert(filePath string) error {
	// Check if the provided file has a .webp extension.
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

	// Generate the output filenames for JPG and PNG.
	base := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	jpgFile := base + ".jpg"
	pngFile := base + ".png"

	// Create and encode the JPEG version.
	jf, err := os.Create(jpgFile)
	if err != nil {
		return fmt.Errorf("unable to create jpg file: %w", err)
	}
	if err = jpeg.Encode(jf, img, &jpeg.Options{Quality: 100}); err != nil {
		jf.Close()
		return fmt.Errorf("failed to encode jpg: %w", err)
	}
	jf.Close()

	// Create and encode the PNG version.
	pf, err := os.Create(pngFile)
	if err != nil {
		return fmt.Errorf("unable to create png file: %w", err)
	}
	if err = png.Encode(pf, img); err != nil {
		pf.Close()
		return fmt.Errorf("failed to encode png: %w", err)
	}
	pf.Close()

	// Get file sizes for the original WebP, JPG, and PNG files.
	webpInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file stats for webp: %w", err)
	}
	jpgInfo, err := os.Stat(jpgFile)
	if err != nil {
		return fmt.Errorf("failed to get file stats for jpg: %w", err)
	}
	pngInfo, err := os.Stat(pngFile)
	if err != nil {
		return fmt.Errorf("failed to get file stats for png: %w", err)
	}

	// Print the file sizes for all three images.
	fmt.Printf("%s size: %d bytes\n", filepath.Base(filePath), webpInfo.Size())
	fmt.Printf("%s size: %d bytes\n", filepath.Base(jpgFile), jpgInfo.Size())
	fmt.Printf("%s size: %d bytes\n", filepath.Base(pngFile), pngInfo.Size())

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
