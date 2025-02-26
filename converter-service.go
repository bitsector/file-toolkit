package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const uploadPath = "./uploads"
const outputPath = "./converted"

func init() {
	// Create directories if they do not exist.
	os.Mkdir(uploadPath, 0755)
	os.Mkdir(outputPath, 0755)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with a maximum memory of 10MB.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the uploaded file.
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check that the file extension is ".webp".
	if strings.ToLower(filepath.Ext(header.Filename)) != ".webp" {
		http.Error(w, "Invalid file type. Only WebP files accepted", http.StatusBadRequest)
		return
	}

	// Decode the WebP image.
	img, err := webp.Decode(file)
	if err != nil {
		http.Error(w, "Error decoding WebP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create output file name and prepare the file path.
	outputName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)) + ".jpg"
	outputFilePath := filepath.Join(outputPath, outputName)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		http.Error(w, "Error creating output file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Encode the image as JPEG.
	if err := jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 95}); err != nil {
		http.Error(w, "Error converting to JPEG: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully converted to: %s", outputName)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	// Create an HTTP/2 server instance.
	h2s := &http2.Server{}

	// Enable h2c (HTTP/2 without TLS) by wrapping the mux handler.
	server := &http.Server{
		Addr:    ":3000",
		Handler: h2c.NewHandler(mux, h2s),
	}

	log.Println("Starting HTTP/2 server on localhost:3000")
	log.Fatal(server.ListenAndServe())
}
