package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// getEnv retrieves environment variables with fallback logic:
// 1. OS environment variables (highest priority)
// 2. .env file values (if loaded)
// 3. Default values (fallback)
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves integer environment variables with fallback logic
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Warn().Str("key", key).Str("value", value).Msg("Invalid integer value in environment, using default")
	}
	return defaultValue
}

// getEnvInt64 retrieves int64 environment variables with fallback logic
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
		log.Warn().Str("key", key).Str("value", value).Msg("Invalid int64 value in environment, using default")
	}
	return defaultValue
}

// getEnvDuration retrieves duration environment variables with fallback logic
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Warn().Str("key", key).Str("value", value).Msg("Invalid duration value in environment, using default")
	}
	return defaultValue
}

func init() {
	// Load environment variables from .env file if present
	// OS environment variables will override .env values
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using OS environment variables and defaults")
	}

	// Configure logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Set configuration variables from environment
	uploadPath = getEnv("UPLOAD_PATH", "./uploads")
	numWorkers = getEnvInt("NUM_WORKERS", 5)
	jobTimeout = getEnvDuration("JOB_TIMEOUT", 30*time.Second)
	workerResultTimeout = getEnvDuration("WORKER_RESULT_TIMEOUT", 1*time.Second)
	jobQueueTimeout = getEnvDuration("JOB_QUEUE_TIMEOUT", 100*time.Millisecond)
	jobQueueSize = getEnvInt("JOB_QUEUE_SIZE", 100)

	// Initialize job queue channel with configured size
	jobQueue = make(chan ConversionJob, jobQueueSize)

	// Create upload directory
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create upload directory")
	}

	// Parse buffer size from environment (default: 10MB)
	bufferSize = getEnvInt64("BUFFER_SIZE", 10<<20)

	// Parse server port from environment (default: 3000)
	port := getEnv("PORT", "3000")
	if port[0] != ':' {
		serverPort = ":" + port
	} else {
		serverPort = port
	}
}
