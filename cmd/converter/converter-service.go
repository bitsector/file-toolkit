package main

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/image/webp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	uploadPath = "./uploads"
	numWorkers = 5
	jobTimeout = 30 * time.Second
)

var (
	serverPort string
	bufferSize int64
	jobQueue   = make(chan ConversionJob, 100) // Buffered job channel
	workerWg   sync.WaitGroup                  // Worker synchronization
)

// ConversionJob represents a WebP conversion task
type ConversionJob struct {
	ID       string
	File     multipart.File
	Header   *multipart.FileHeader
	Result   chan ConversionResult
	Deadline time.Time
}

// ConversionResult holds the conversion output or error
type ConversionResult struct {
	Buffer     bytes.Buffer
	OutputName string
	Error      error
	Metrics    ConversionMetrics
}

// ConversionMetrics holds performance and size metrics
type ConversionMetrics struct {
	ConversionDuration time.Duration
	OriginalSizeMB     float64
	ConvertedSizeMB    float64
	OriginalBytes      int64
	ConvertedBytes     int
}

// worker processes conversion jobs from the job queue
func worker(ctx context.Context) {
	defer workerWg.Done()
	for {
		select {
		case job := <-jobQueue:
			result := processConversionJob(job)
			select {
			case job.Result <- result:
				// Result sent successfully
			case <-time.After(1 * time.Second):
				// Timeout sending result, client may have disconnected
				log.Warn().Str("job_id", job.ID).Msg("Failed to send result - client timeout")
			}
		case <-ctx.Done():
			log.Info().Msg("Worker shutting down")
			return
		}
	}
}

// processConversionJob handles the actual WebP to JPEG conversion
func processConversionJob(job ConversionJob) ConversionResult {
	startTime := time.Now()

	// Check if job has expired
	if time.Now().After(job.Deadline) {
		return ConversionResult{
			Error: fmt.Errorf("job expired"),
		}
	}

	// Validate file extension
	if strings.ToLower(filepath.Ext(job.Header.Filename)) != ".webp" {
		return ConversionResult{
			Error: fmt.Errorf("invalid file type. Only WebP files accepted"),
		}
	}

	// Decode WebP image
	img, err := webp.Decode(job.File)
	if err != nil {
		return ConversionResult{
			Error: fmt.Errorf("error decoding WebP: %w", err),
		}
	}

	// Prepare output filename
	outputName := strings.TrimSuffix(job.Header.Filename, filepath.Ext(job.Header.Filename)) + ".jpg"

	// Encode to JPEG
	var imgBuffer bytes.Buffer
	if err := jpeg.Encode(&imgBuffer, img, &jpeg.Options{Quality: 95}); err != nil {
		return ConversionResult{
			Error: fmt.Errorf("error converting to JPEG: %w", err),
		}
	}

	// Calculate metrics
	conversionDuration := time.Since(startTime)
	originalSizeMB := float64(job.Header.Size) / (1024 * 1024)
	convertedSizeMB := float64(imgBuffer.Len()) / (1024 * 1024)

	return ConversionResult{
		Buffer:     imgBuffer,
		OutputName: outputName,
		Error:      nil,
		Metrics: ConversionMetrics{
			ConversionDuration: conversionDuration,
			OriginalSizeMB:     originalSizeMB,
			ConvertedSizeMB:    convertedSizeMB,
			OriginalBytes:      job.Header.Size,
			ConvertedBytes:     imgBuffer.Len(),
		},
	}
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(bufferSize); err != nil {
		log.Error().Err(err).Msg("Error parsing form")
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving file")
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create result channel
	resultChan := make(chan ConversionResult, 1)
	jobID := fmt.Sprintf("conversion-%d", time.Now().UnixNano())

	// Create conversion job
	job := ConversionJob{
		ID:       jobID,
		File:     file,
		Header:   header,
		Result:   resultChan,
		Deadline: time.Now().Add(jobTimeout),
	}

	// Submit job to worker pool with timeout
	select {
	case jobQueue <- job:
		log.Info().Str("job_id", jobID).Str("filename", header.Filename).Msg("Job submitted to worker pool")
	case <-time.After(100 * time.Millisecond):
		log.Warn().Str("job_id", jobID).Msg("Job queue full - server busy")
		http.Error(w, "Server busy, please try again later", http.StatusServiceUnavailable)
		return
	}

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.Error != nil {
			log.Error().Err(result.Error).Str("job_id", jobID).Msg("Conversion failed")

			// Determine appropriate HTTP status code based on error
			statusCode := http.StatusInternalServerError
			if strings.Contains(result.Error.Error(), "invalid file type") {
				statusCode = http.StatusBadRequest
			} else if strings.Contains(result.Error.Error(), "job expired") {
				statusCode = http.StatusRequestTimeout
			}

			http.Error(w, result.Error.Error(), statusCode)
			return
		}

		// Log conversion metrics
		log.Info().
			Str("job_id", jobID).
			Str("original_size", fmt.Sprintf("%.2f MB (%d bytes)", result.Metrics.OriginalSizeMB, result.Metrics.OriginalBytes)).
			Str("converted_size", fmt.Sprintf("%.2f MB (%d bytes)", result.Metrics.ConvertedSizeMB, result.Metrics.ConvertedBytes)).
			Str("conversion_time", result.Metrics.ConversionDuration.String()).
			Msg("File conversion completed")

		// Send response
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "attachment; filename="+result.OutputName)
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(result.Buffer.Bytes()); err != nil {
			log.Error().Err(err).Str("job_id", jobID).Msg("Error writing response")
		}

	case <-time.After(jobTimeout):
		log.Error().Str("job_id", jobID).Msg("Conversion timeout")
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	}

	close(resultChan)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	log.Info().Int("num_workers", numWorkers).Msg("Starting worker pool")
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go worker(ctx)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	// Setup HTTP/2 server
	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    serverPort,
		Handler: h2c.NewHandler(mux, h2s),
	}

	log.Info().Msgf("Starting HTTP/2 server on %s", serverPort)

	// Start server and handle graceful shutdown
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server failed")
	}

	// Signal workers to shutdown
	cancel()

	// Wait for all workers to finish
	log.Info().Msg("Waiting for workers to finish...")
	workerWg.Wait()
	log.Info().Msg("Server shutdown complete")
}
