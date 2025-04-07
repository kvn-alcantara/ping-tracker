package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
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

var (
	stdout      io.Writer = os.Stdout
	stdoutMutex sync.RWMutex
)

func SetOutput(w io.Writer) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	stdout = w
}

func (d *TerminalDisplay) ClearScreen() {
	stdoutMutex.RLock()
	defer stdoutMutex.RUnlock()
	fmt.Fprint(stdout, "\033[H\033[2J")
}

func (d *TerminalDisplay) PrintHeader(title string) {
	stdoutMutex.RLock()
	defer stdoutMutex.RUnlock()
	fmt.Fprintln(stdout, d.cyan(title))
	fmt.Fprintln(stdout, "-------------------------------------------------------")
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

	stdoutMutex.RLock()
	defer stdoutMutex.RUnlock()
	fmt.Fprintf(stdout, "%-30s [%s]%s\n", url, statusColor(status), latencyStr)
}

var _ Display = (*TerminalDisplay)(nil)
