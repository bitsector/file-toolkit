package main

import (
	"bytes"
	"image/jpeg"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/webp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	uploadPath = "./uploads"
)

var (
	serverPort string
	bufferSize int64
)

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(bufferSize); err != nil {
		log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving file")
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if strings.ToLower(filepath.Ext(header.Filename)) != ".webp" {
		log.Warn().Str("filename", header.Filename).Msg("Invalid file type")
		http.Error(w, "Invalid file type. Only WebP files accepted", http.StatusBadRequest)
		return
	}

	img, err := webp.Decode(file)
	if err != nil {
		log.Error().Err(err).Msg("Error decoding WebP")
		http.Error(w, "Error decoding WebP: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outputName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename)) + ".jpg"

	var imgBuffer bytes.Buffer
	if err := jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 95}); err != nil {
		log.Error().Err(err).Msg("Error converting to JPEG")
		http.Error(w, "Error converting to JPEG: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename="+outputName)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(imgBuffer.Bytes()); err != nil {
		log.Error().Err(err).Msg("Error writing response")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    serverPort,
		Handler: h2c.NewHandler(mux, h2s),
	}

	log.Info().Msgf("Starting HTTP/2 server on %s", serverPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
