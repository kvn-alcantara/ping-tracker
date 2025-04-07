package usecase

import (
	"sync"

	"github.com/kvn-alcantara/ping-tracker/internal/delivery/cli"
	"github.com/kvn-alcantara/ping-tracker/internal/entity"
	"github.com/kvn-alcantara/ping-tracker/internal/repository"
)

type PingTracker struct {
	pinger     repository.Pinger
	display    cli.Display
	logger     cli.Logger
	urls       []string
	hostStatus map[string]entity.PingResult
	updateChan chan entity.PingResult
	wg         sync.WaitGroup
}
