package factory

import (
	"runtime"

	"github.com/kvn-alcantara/ping-tracker/internal/repository"
)

func NewPinger() repository.Pinger {
	if runtime.GOOS == "linux" {
		return repository.NewHttpPinger()
	}
	return repository.NewProBingPinger()
}
