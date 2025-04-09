package cli

import (
	"log"
	"os"
)

// Logger is an interface for logging messages.
type Logger interface {
	// Log writes a message to the log.
	Log(message string)
}

// FileLogger is a logger that writes logs to a file.
type FileLogger struct {
	file   *os.File
	logger *log.Logger
}

// NewFileLogger creates a new FileLogger that writes to the specified file path.
// If the file does not exist, it will be created. Logs are appended to the file.
func NewFileLogger(filePath string) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{
		file:   file,
		logger: log.New(file, "[HOST-MONITOR] ", log.LstdFlags),
	}, nil
}

// Log writes a message to the log file.
func (l *FileLogger) Log(message string) {
	l.logger.Println(message)
}

// Close closes the log file.
func (l *FileLogger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

var _ Logger = (*FileLogger)(nil)
