package usecase

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kvn-alcantara/ping-tracker/internal/delivery/cli"
	"github.com/kvn-alcantara/ping-tracker/internal/entity"
	"github.com/kvn-alcantara/ping-tracker/internal/repository"
)

type Monitor struct {
	pinger     repository.Pinger
	display    cli.Display
	logger     cli.Logger
	urls       []string
	hostStatus map[string]entity.Host
	updateChan chan entity.Host
	wg         sync.WaitGroup
}

func NewMonitor(pinger repository.Pinger, display cli.Display, logger cli.Logger, urls []string) *Monitor {
	return &Monitor{
		pinger:     pinger,
		display:    display,
		logger:     logger,
		urls:       urls,
		hostStatus: make(map[string]entity.Host),
		updateChan: make(chan entity.Host, 100),
	}
}

func (m *Monitor) StartMonitoring() {
	for _, url := range m.urls {
		m.hostStatus[url] = entity.Host{URL: url, Status: "Resolving..."}
		m.wg.Add(1)
		go m.monitorContinuously(url)
	}
}

func (m *Monitor) monitorContinuously(url string) {
	defer m.wg.Done()

	resolver := net.Resolver{}

	for {
		ips, err := resolver.LookupIP(context.Background(), "ip", url)
		if err != nil {
			m.logger.Log(fmt.Sprintf("DNS Error for %s: %v", url, err))
			m.updateChan <- entity.Host{URL: url, Status: "DNS Error", Error: err}
			time.Sleep(1 * time.Second)
			continue
		}
		if len(ips) == 0 {
			m.logger.Log(fmt.Sprintf("No IP Found for %s", url))
			m.updateChan <- entity.Host{URL: url, Status: "No IP Found"}
			time.Sleep(1 * time.Second)
			continue
		}
		ip := ips[0].String()
		latency, err := m.pinger.Ping(ip)
		if err != nil {
			m.logger.Log(fmt.Sprintf("Ping failed for %s (%s): %v", url, ip, err))
			m.updateChan <- entity.Host{URL: url, Status: "Offline", Error: err}
		} else {
			m.updateChan <- entity.Host{URL: url, Status: "Online", Latency: latency}
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *Monitor) Run(done <-chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	m.display.PrintHeader("Real-time Host Monitor - Press Ctrl+C to Exit")

	for {
		select {
		case result := <-m.updateChan:
			m.hostStatus[result.URL] = result
		case <-ticker.C:
			m.display.ClearScreen()
			m.display.PrintHeader("Real-time Host Monitor - Press Ctrl+C to Exit")
			for _, url := range m.urls {
				status := m.hostStatus[url].Status
				latency := m.hostStatus[url].Latency
				m.display.PrintStatus(url, status, latency)
			}
		case <-done:
			return
		}
	}
	m.wg.Wait()
}
