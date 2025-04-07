package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileLogger(t *testing.T) {
	withTempFile(t, func(tempFile string) {
		logger, err := NewFileLogger(tempFile)
		assert.NoError(t, err, "expected no error when creating a new FileLogger")
		assert.NotNil(t, logger, "expected logger to be initialized")

		_, err = os.Stat(tempFile)
		assert.NoError(t, err, "expected log file to be created")
	})
}

func TestFileLoggerLog(t *testing.T) {
	withTempFile(t, func(tempFile string) {
		logger, err := NewFileLogger(tempFile)
		assert.NoError(t, err, "expected no error when creating a new FileLogger")
		defer logger.Close()

		testMessage := "This is a test log message"
		logger.Log(testMessage)

		content, err := os.ReadFile(tempFile)
		assert.NoError(t, err, "expected no error when reading the log file")
		assert.Contains(t, string(content), testMessage, "expected log file to contain the test message")
	})
}

func TestFileLoggerClose(t *testing.T) {
	withTempFile(t, func(tempFile string) {
		logger, err := NewFileLogger(tempFile)
		assert.NoError(t, err, "expected no error when creating a new FileLogger")

		logger.Close()

		_, err = os.Stat(tempFile)
		assert.NoError(t, err, "expected log file to still exist after closing the logger")
	})
}

func withTempFile(t *testing.T, testFunc func(tempFile string)) {
	tempDir := t.TempDir()
	tempFile := tempDir + "/test_logger.log"
	testFunc(tempFile)
}
