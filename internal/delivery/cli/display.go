package cli

import (
	"fmt"
	"time"

	"github.com/mgutz/ansi"
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

func NewTerminalDisplay() *TerminalDisplay {
	return &TerminalDisplay{
		green:  ansi.ColorFunc("green+"),
		red:    ansi.ColorFunc("red+"),
		yellow: ansi.ColorFunc("yellow+"),
		cyan:   ansi.ColorFunc("cyan+"),
	}
}

func (d *TerminalDisplay) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (d *TerminalDisplay) PrintHeader(title string) {
	fmt.Println(d.cyan(title))
	fmt.Println("-------------------------------------------------------")
}

func (d *TerminalDisplay) PrintStatus(url string, status string, latency time.Duration) {
	var statusColor func(string) string

	switch status {
	case "Online":
		statusColor = d.green
	case "Offline", "DNS Error", "Ping Setup Error", "No IP Found":
		statusColor = d.red
	case "Resolving...", "Pinging...":
		statusColor = d.yellow
	default:
		statusColor = ansi.ColorFunc("")
	}

	latencyStr := ""
	if status == "Online" && latency > 0 {
		latencyStr = fmt.Sprintf(" (Latency: %s)", latency.Round(time.Millisecond))
	}

	fmt.Printf("%-30s [%s]%s\n", url, statusColor(status), latencyStr)
}

var _ Display = (*TerminalDisplay)(nil)
