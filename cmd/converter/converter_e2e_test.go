package main

import (
	"bytes"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestWebpToJpgE2E(t *testing.T) {
	// Open the sample webp file
	file, err := os.Open("../../samples/meme.webp")
	if err != nil {
		t.Fatalf("failed to open sample webp file: %v", err)
	}
	defer file.Close()

	// Prepare multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "meme.webp")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		t.Fatalf("failed to copy file data: %v", err)
	}
	writer.Close()

	// Send POST request to the converter
	resp, err := http.Post("http://localhost:3000/convert", writer.FormDataContentType(), &body)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Read response body
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Assert that the response is a non-trivial JPEG file
	if len(respBytes) < 1000 {
		t.Fatalf("response too small to be a valid JPEG: got %d bytes", len(respBytes))
	}

	// Try to decode as JPEG
	_, err = jpeg.Decode(bytes.NewReader(respBytes))
	if err != nil {
		t.Fatalf("response is not a valid JPEG: %v", err)
	}
}
