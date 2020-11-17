package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "url"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
