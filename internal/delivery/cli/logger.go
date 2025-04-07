package cli

import (
	"log"
	"os"
)

type Logger interface {
	Log(message string)
}

type FileLogger struct {
	logFile *os.File
	logger  *log.Logger
}
