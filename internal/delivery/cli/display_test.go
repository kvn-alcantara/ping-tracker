package cli

import (
	"bytes"
	"testing"
	"time"
)

func TestTerminalDisplay(t *testing.T) {
	display := NewTerminalDisplay()

	t.Run("PrintHeader", func(t *testing.T) {
		var buf bytes.Buffer
		// Temporarily redirect stdout to our buffer
		old := stdout
		stdout = &buf
		defer func() { stdout = old }()

		display.PrintHeader("Test Title")
		output := buf.String()
		if output == "" {
			t.Error("PrintHeader should output text")
		}
	})

	t.Run("PrintStatus", func(t *testing.T) {
		var buf bytes.Buffer
		old := stdout
		stdout = &buf
		defer func() { stdout = old }()

		tests := []struct {
			url     string
			status  string
			latency time.Duration
		}{
			{"example.com", "Online", 100 * time.Millisecond},
			{"example.com", "Offline", 0},
			{"example.com", "Resolving...", 0},
		}

		for _, tt := range tests {
			buf.Reset()
			display.PrintStatus(tt.url, tt.status, tt.latency)
			output := buf.String()
			if output == "" {
				t.Errorf("PrintStatus should output text for status %s", tt.status)
			}
		}
	})

	t.Run("ClearScreen", func(t *testing.T) {
		var buf bytes.Buffer
		old := stdout
		stdout = &buf
		defer func() { stdout = old }()

		display.ClearScreen()
		output := buf.String()
		if output == "" {
			t.Error("ClearScreen should output ANSI escape codes")
		}
	})
}
