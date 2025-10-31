package service

import (
	"testing"

	"go.uber.org/zap/zapcore"
	"quizizz.com/internal/logger"
)

// DisableLoggingForBenchmark temporarily disables logging for benchmarks
func DisableLoggingForBenchmark(b *testing.B) {
	// Save the current level
	prevLevel := logger.GetLevel()

	// Set log level to a level that suppresses all output
	logger.SetLevel(zapcore.FatalLevel + 1)

	// Restore the original level when the benchmark is done
	b.Cleanup(func() {
		logger.SetLevel(prevLevel)
	})
}
