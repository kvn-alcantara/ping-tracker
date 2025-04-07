package cli

import (
	"log"
	"os"
)

type Logger interface {
	Log(message string)
}

type FileLogger struct {
	file   *os.File
	logger *log.Logger
}

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

func (l *FileLogger) Log(message string) {
	l.logger.Println(message)
}

func (l *FileLogger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

var _ Logger = (*FileLogger)(nil)
