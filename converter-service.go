package main

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/webp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	uploadPath = "./uploads"
	outputPath = "./converted"
)

var (
	serverPort string
	bufferSize int64
)

func init() {
	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using defaults")
	}

	// Set the zerolog time format and global logging level.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create directories if they do not exist.
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create upload directory")
	}
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create output directory")
	}

	// Read buffer size from environment variable, default to 10MB if not set.
	bufferSizeStr := os.Getenv("BUFFER_SIZE")
	if bufferSizeStr == "" {
		bufferSize = 10 << 20 // 10MB default
	} else {
		val, err := strconv.ParseInt(bufferSizeStr, 10, 64)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid BUFFER_SIZE value, using default (10MB)")
			bufferSize = 10 << 20
		} else {
			bufferSize = val
		}
	}

	// Set server port, default to 3000 if not provided.
	port := os.Getenv("PORT")
	if port == "" {
		serverPort = ":3000"
	} else {
		// Prepend ':' if not present.
		if port[0] != ':' {
			serverPort = ":" + port
		} else {
			serverPort = port
		}
	}
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with the bufferSize specified.
	if err := r.ParseMultipartForm(bufferSize); err != nil {
		log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the uploaded file.
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving file")
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check that the file extension is ".webp".
	if strings.ToLower(filepath.Ext(header.Filename)) != ".webp" {
		log.Warn().Str("filename", header.Filename).Msg("Invalid file type")
		http.Error(w, "Invalid file type. Only WebP files accepted", http.StatusBadRequest)
		return
	}

	// Decode the WebP image.
	img, err := webp.Decode(file)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding WebP")
		http.Error(w, "Error decoding WebP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create output file name and prepare the file path.
	outputName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)) + ".jpg"
	outputFilePath := filepath.Join(outputPath, outputName)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Error creating output file")
		http.Error(w, "Error creating output file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Encode the image as JPEG.
	if err := jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 95}); err != nil {
		log.Error().Err(err).Msg("Error converting to JPEG")
		http.Error(w, "Error converting to JPEG: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Str("outputFile", outputFilePath).Msg("Successfully converted file")
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "Successfully converted to: %s", outputName); err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	// Create an HTTP/2 server instance.
	h2s := &http2.Server{}

	// Enable h2c (HTTP/2 without TLS) by wrapping the mux handler.
	server := &http.Server{
		Addr:    serverPort,
		Handler: h2c.NewHandler(mux, h2s),
	}

	log.Info().Msgf("Starting HTTP/2 server on %s", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
