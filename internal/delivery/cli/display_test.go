package cli

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestTerminalDisplay(t *testing.T) {
	display := NewTerminalDisplay()

	t.Run("PrintHeader", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(os.Stdout)

		display.PrintHeader("Test Title")
		output := buf.String()
		if output == "" {
			t.Error("PrintHeader should output text")
		}
	})

	t.Run("PrintStatus", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		defer SetOutput(os.Stdout)

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
		SetOutput(&buf)
		defer SetOutput(os.Stdout)

		display.ClearScreen()
		output := buf.String()
		if output == "" {
			t.Error("ClearScreen should output ANSI escape codes")
		}
	})
}
