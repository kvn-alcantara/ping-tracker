package cli

import (
	"time"
)

type Display interface {
	ClearScreen()
	PrintHeader(title string)
	PrintStatus(url string, status string, latency time.Duration)
}

type TerminalDisplay struct {
	green  func(string) string
	red    func(string) string
	yellow func(string) string
	cyan   func(string) string
}
