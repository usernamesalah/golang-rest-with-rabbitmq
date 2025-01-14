package logger

import (
	"os"
	"tenant/infrastructure/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger initializes and returns a Logrus logger with file rotation.
func NewLogger(cfg config.Config) *logrus.Logger {
	// Create a new instance of Logrus logger
	logger := logrus.New()

	// Configure Logrus to use Lumberjack for log rotation
	logger.SetOutput(&lumberjack.Logger{
		Filename:   cfg.Logging.File,
		MaxSize:    cfg.Logging.MaxSize,    // Maximum size in MB
		MaxBackups: cfg.Logging.MaxBackups, // Maximum number of old log files to retain
		MaxAge:     cfg.Logging.MaxAge,     // Maximum age in days to retain old log files
	})

	// Set the log format to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level (default to Info)
	logger.SetLevel(logrus.InfoLevel)

	// Also log to the console (optional)
	logger.SetOutput(os.Stdout)

	return logger
}
