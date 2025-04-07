package main

import (
	"fmt"
	"os"

	"github.com/kvn-alcantara/ping-tracker/internal/delivery/cli"
	"github.com/kvn-alcantara/ping-tracker/internal/repository"
)

func main() {
	logger, err := cli.NewFileLogger("host_monitor.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	cli.NewTerminalDisplay()
	repository.NewProBingPinger()

	logger.Log("Ping Tracker started")
}
