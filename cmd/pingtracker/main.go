package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kvn-alcantara/ping-tracker/internal/delivery/cli"
	"github.com/kvn-alcantara/ping-tracker/internal/repository"
	"github.com/kvn-alcantara/ping-tracker/internal/usecase"
)

func main() {
	urls := []string{"google.com", "github.com", "stackoverflow.com", "reddit.com"}

	logger, err := cli.NewFileLogger("host_monitor.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	display := cli.NewTerminalDisplay()
	pinger := repository.NewProBingPinger()
	monitor := usecase.NewMonitor(pinger, display, logger, urls)

	monitor.StartMonitoring()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		close(done)
	}()

	monitor.Run(done)
}
