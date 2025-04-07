package main

import (
	"github.com/kvn-alcantara/ping-tracker/internal/delivery/cli"
)

func main() {
	display := cli.NewTerminalDisplay()

	display.PrintHeader("Ping Tracker")
}
