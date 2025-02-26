package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using defaults")
	}

	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create upload directory
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create upload directory")
	}

	// Parse buffer size from environment
	bufferSizeStr := os.Getenv("BUFFER_SIZE")
	if bufferSizeStr == "" {
		bufferSize = 10 << 20 // Default 10MB
	} else {
		val, err := strconv.ParseInt(bufferSizeStr, 10, 64)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid BUFFER_SIZE value, using default")
			bufferSize = 10 << 20
		} else {
			bufferSize = val
		}
	}

	// Parse server port from environment
	port := os.Getenv("PORT")
	if port == "" {
		serverPort = ":3000"
	} else {
		if port[0] != ':' {
			serverPort = ":" + port
		} else {
			serverPort = port
		}
	}
}
